use env_logger;
use log::{info, error};
use robot_hal::{Motor, Sensor, N20Motor, DummySensor};

fn main() {
    // Initialize logger (requirement 11)
    env_logger::init();
    
    info!("Robot HAL Starting...");
    
    // Print hello from main
    println!("Hello from Robot Main!");
    
    // Try to create a motor instance - handle errors without panicking (requirement 7)
    match N20Motor::new(18, 23, 24) {
        Ok(mut motor) => {
            if let Err(e) = motor.say_hello() {
                error!("Failed to say hello from motor: {}", e);
            }
        }
        Err(e) => {
            error!("Failed to create motor: {}", e);
        }
    }
    
    // Create and test dummy sensor
    let mut sensor = DummySensor::new();
    if let Err(e) = sensor.say_hello() {
        error!("Failed to say hello from sensor: {}", e);
    }
    
    info!("Robot HAL Shutdown Complete");
} 