package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client mqtt.Client
)

func StartMQTTClient(broker *string, clientID *string) {
	// Create a new MQTT client
	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID(*clientID)

	// Set up the connection lost handler
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
	}

	// Set up the message handler
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to broker")
	}

	// Create and start the MQTT client
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Wait for interrupt signal to gracefully shut down
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, os.Interrupt)
	//<-sigChan
}

func DisconnectMQTTClient() {
	// Clean up
	client.Disconnect(250)
	fmt.Println("Client disconnected")
}

// publishMessages sends a message to the MQTT broker
func PublishMessages(topic string, message string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()

	if token.Error() != nil {
		// Handle the publish error
		fmt.Printf("Error publishing message: %v\n", token.Error())
	} else {
		fmt.Printf("Published message: %s\n", message)
	}
}
