package mqtt

import (
	"encoding/json"
	"telemetry/include/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	QoSAtMostOnce  = 0
	QoSAtLeastOnce = 1
	QoSExactlyOnce = 2
)

type MQTTTelemetryClient struct {
	client    mqtt.Client
	robotID   string
	brokerURL string
	log       *logger.Logger
}

func NewMQTTTelemetryClient(robotID, brokerURL string) *MQTTTelemetryClient {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(robotID).
		SetAutoReconnect(true).
		SetCleanSession(false).
		SetOrderMatters(false)

	client := mqtt.NewClient(opts)

	return &MQTTTelemetryClient{
		client:    client,
		robotID:   robotID,
		brokerURL: brokerURL,
		log:       logger.New(logger.INFO),
	}
}

func (m *MQTTTelemetryClient) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(m.brokerURL).
		SetClientID(m.robotID).
		SetAutoReconnect(true).
		SetCleanSession(false).
		SetOrderMatters(false)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	m.client = client
	return nil
}

func (m *MQTTTelemetryClient) SubscribeToCommands(callback mqtt.MessageHandler) error {
	topic := "robots/" + m.robotID + "/commands"
	token := m.client.Subscribe(topic, QoSAtLeastOnce, callback)
	token.Wait()
	return token.Error()
}

func (m *MQTTTelemetryClient) PublishTelemetry(messageType string, data interface{}) error {
	topic := "robots/" + m.robotID + "/telemetry/" + messageType

	jsonData, err := json.Marshal(data)
	if err != nil {
		m.log.Error("failed to marshal telemetry data: %v", err)
		return err
	}

	token := m.client.Publish(topic, QoSAtLeastOnce, false, jsonData)
	token.Wait()
	return token.Error()
}

func (m *MQTTTelemetryClient) BrokerURL() string {
	return m.brokerURL
}

func (m *MQTTTelemetryClient) IsConnected() bool {
	return m.client != nil && m.client.IsConnected()
}

// ... move all the MQTTTelemetryClient methods here ...
