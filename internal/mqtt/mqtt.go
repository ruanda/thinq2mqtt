package mqtt

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Server string
	Topic  string
}

type Client struct {
	config Config
	client MQTT.Client
}

func NewClient(config Config) *Client {
	clientOpts := MQTT.NewClientOptions()
	clientOpts.AddBroker(config.Server)

	c := &Client{
		config: config,
		client: MQTT.NewClient(clientOpts),
	}

	return c
}

func (c *Client) Connect(timeout time.Duration) error {
	token := c.client.Connect()
	if !token.WaitTimeout(timeout) {
		return "TODO"
	}
	return token.Error()
}
