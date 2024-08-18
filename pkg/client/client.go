// client.go
package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	// Parametro per il file di configurazione kubeconfig
	kubeconfig := flag.String("kubeconfig", "/root/.kube/config", "absolute path to the kubeconfig file")
	// Parametro per l'URL del broker MQTT
	broker := flag.String("broker", "tcp://broker.hivemq.com:1883", "MQTT broker URL")
	// Parametro per il topic MQTT
	topic := flag.String("topic", "test/topic123679", "MQTT topic")
	// Parametro per l'ID del client MQTT
	clientID := flag.String("clientID", "clientId-KI02s7qxUZ", "MQTT client ID")
	// Parametro per la porta del server Latency Meter
	latencyMeterServerPort := flag.String("lmport", "8080", "Latency Meter server port")
	// Parametro per il selettore di etichette per filtrare i pod
	labelSelector := flag.String("label", "feature=latency-aware-deployment", "label selector to filter pods")
	// Parametro per il namespace
	namespace := flag.String("namespace", "", "namespace to filter pods")
	// Parametro per il namespace
	measurementsArg := flag.String("measurements", "10", "number of measurements to take")

	flag.Parse()

	// Converti il numero di misurazioni in un intero
	measurements, err := strconv.Atoi(*measurementsArg)
	if err != nil {
		fmt.Println("Error converting port to integer:", err)
		return
	}

	// Connessione al cluster Kubernetes
	nodes, err := GetNodeIPs(kubeconfig, labelSelector, namespace)
	if err != nil {
		fmt.Println("Error getting node IPs:", err)
		return
	}

	// Connessione al server MQTT
	StartMQTTClient(broker, clientID)
	defer DisconnectMQTTClient()

	for _, node := range nodes {
		fmt.Printf("Node: %s\n", node.Name)
		fmt.Printf("Pod: %s\n", node.PodName)
		fmt.Printf("Namespace: %s\n", node.Namespace)
		fmt.Printf("Internal IP: %s\n", node.InternalIP)
		fmt.Printf("External IP: %s\n", node.ExternalIP)

		// Connessione al server Latency Meter
		address := node.ExternalIP
		if address == "" {
			address = node.InternalIP
		}
		address += ":" + *latencyMeterServerPort
		fmt.Printf("Connecting to Latency Meter server at %s\n", address)

		// Connessione al server UDP sulla porta 8080
		port, err := strconv.Atoi(*latencyMeterServerPort)
		if err != nil {
			fmt.Println("Error converting port to integer:", err)
			continue
		}

		serverAddr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP(node.InternalIP),
		}

		// Connect to the server
		conn, err := net.DialUDP("udp", nil, &serverAddr)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			continue
		}
		defer conn.Close()

		// Repeat the process 10 times and measure the average RTT
		successfulMeasurements := 0
		totalRTT := time.Duration(0)
		for i := 0; i < measurements; i++ {
			// Send a request
			startTime := time.Now()
			_, err = conn.Write([]byte("ping"))
			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}

			// Read the response
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading response:", err)
				continue
			}

			rtt := time.Since(startTime)
			fmt.Printf("RTT: %v\n", rtt)

			successfulMeasurements++
			totalRTT += rtt
		}

		if successfulMeasurements == 0 {
			continue
		}

		averageRTT := totalRTT / time.Duration(successfulMeasurements)
		fmt.Printf("Average RTT: %v\n", averageRTT)

		PublishMessages(*topic, fmt.Sprintf("Pod: %s\nNamespace: %s\nRTT: %s\nTimestamp: %s", node.PodName, node.Namespace, averageRTT.String(), time.Now().Format(time.RFC3339)))
	}

}
