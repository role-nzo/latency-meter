// server.go
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	// Parametro per la porta del server Latency Meter
	latencyMeterServerPort := flag.String("lmport", "8080", "Latency Meter server port")
	flag.Parse()

	// Avvia il server UDP
	port, _ := strconv.Atoi(*latencyMeterServerPort)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("localhost"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Errore nell'aprire la porta UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Server UDP in ascolto su porta " + *latencyMeterServerPort + "...")

	buffer := make([]byte, 4) // Buffer per leggere il messaggio

	for {
		// Leggi il messaggio inviato dal client
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Errore nella lettura del messaggio:", err)
			continue
		}

		// Se il messaggio è "ping", rispondi con "pong"
		if string(buffer) == "ping" {
			_, err = conn.WriteToUDP([]byte("pong"), addr)
			if err != nil {
				fmt.Println("Errore nell'invio della risposta:", err)
			}
		}
	}
}
