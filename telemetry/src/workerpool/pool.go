package workerpool

import (
	"context"
	"fmt"
	"sync"
	"telemetry/include/logger"
	"time"
)

// Job represents a unit of work to be performed
type Job struct {
	Name     string
	Execute  func(ctx context.Context) error
	OnError  func(error)
	Retries  int
	Critical bool // If true, failed critical jobs will trigger system shutdown
}

// WorkerPool manages concurrent job execution with error handling
type WorkerPool struct {
	workers    int
	jobQueue   chan Job
	results    chan error
	log        *logger.Logger
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	errorCount int64
	mu         sync.RWMutex
	done       chan struct{}
}

func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan Job, workers*2),
		results:  make(chan error, workers*2),
		log:      logger.New(logger.INFO),
		ctx:      ctx,
		cancel:   cancel,
		done:     make(chan struct{}),
	}
}

func (wp *WorkerPool) Start() {
	wp.log.Info("Starting worker pool with %d workers", wp.workers)
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// Error monitoring goroutine
	go wp.monitorErrors()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			wp.log.Debug("Worker %d shutting down", id)
			return
		case job := <-wp.jobQueue:
			wp.executeJob(id, job)
		}
	}
}

func (wp *WorkerPool) executeJob(workerID int, job Job) {
	attempt := 0
	for attempt <= job.Retries {
		wp.log.Debug("Worker %d attempting job %s (attempt %d/%d)",
			workerID, job.Name, attempt+1, job.Retries+1)

		err := job.Execute(wp.ctx)
		if err == nil {
			wp.log.Debug("Worker %d completed job %s successfully", workerID, job.Name)
			return
		}

		wp.log.Warn("Worker %d failed job %s (attempt %d/%d): %v",
			workerID, job.Name, attempt+1, job.Retries+1, err)

		if job.OnError != nil {
			job.OnError(err)
		}

		if job.Critical {
			wp.log.Error("Critical job %s failed, initiating shutdown", job.Name)
			wp.results <- err
			return
		}

		attempt++
		if attempt <= job.Retries {
			backoff := time.Second * time.Duration(attempt)
			wp.log.Debug("Worker %d backing off for %v before retrying job %s",
				workerID, backoff, job.Name)
			time.Sleep(backoff)
		}
	}

	wp.mu.Lock()
	wp.errorCount++
	wp.mu.Unlock()
}

func (wp *WorkerPool) monitorErrors() {
	for {
		select {
		case <-wp.ctx.Done():
			return
		case err := <-wp.results:
			if err != nil {
				wp.log.Error("Critical error encountered: %v", err)
				wp.Shutdown()
				return
			}
		}
	}
}

func (wp *WorkerPool) Submit(job Job) {
	select {
	case <-wp.ctx.Done():
		wp.log.Warn("Cannot submit job %s: worker pool is shutting down", job.Name)
		return
	case wp.jobQueue <- job:
		wp.log.Debug("Submitted job: %s", job.Name)
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()

	// Add timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Wait for workers with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		wp.log.Info("Worker pool shut down gracefully")
	case <-shutdownCtx.Done():
		wp.log.Warn("Worker pool shutdown timed out")
	}

	close(wp.jobQueue)
	close(wp.results)
	close(wp.done)
}

func (wp *WorkerPool) GetErrorCount() int64 {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.errorCount
}

func (wp *WorkerPool) Context() context.Context {
	return wp.ctx
}

// Add this method to check if a critical job succeeded
func (wp *WorkerPool) WaitForCriticalJob(jobName string, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	done := make(chan error, 1)
	go func() {
		for {
			select {
			case <-wp.ctx.Done():
				done <- context.Canceled
				return
			case err := <-wp.results:
				done <- err
				return
			}
		}
	}()

	select {
	case err := <-done:
		return err
	case <-timer.C:
		return fmt.Errorf("timeout waiting for critical job: %s", jobName)
	}
}

func (wp *WorkerPool) Done() <-chan struct{} {
	return wp.done
}
