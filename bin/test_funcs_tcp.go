package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

func testPrefixGeneratorLoopTCP(coordinatorAddress string, port int) func() {
	return func() {

		registrationData := []byte{
			0x53, 0x54, // "ST"
			0x00, 0x05, // payload length = 5
			0x01,       // REGISTER
			0x00, 0x01, // TopologyID = 1
			0x00, 0x03, // NodeID = 3
		}

		coordinator_data := make([]byte, 1024)

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", coordinatorAddress, port))

		if err != nil {

			Debug("PrefixGenerator: Can't open the connection to Coordinator  ", "hostname", coordinatorAddress, "port", port)
			log.Fatal("Quittig ..")
		}

		wb, err := conn.Write(registrationData) //write can be in batches
		if err != nil {
			Debug("Cannot Write to Coordinator ", "address", coordinatorAddress)
		}
		Debug("Written to Coordinator ", "bytes", wb)
		rb, err := conn.Read(coordinator_data) //read can be in batches

		if err != nil {

			Debug("Cannot Read From Coordinator ")
		}
		Debug("Read From Coordinator ", "bytes", rb)
		//if all is okay push to coordinator check repsonse
		Debug("data", coordinator_data[0:rb])
		data_frame := []byte{0x53, 0x54, 0, 2, 3, 0}

		for {
			value := rand.IntN(256)
			interval := rand.NormFloat64()*0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			//	timer := time.NewTimer(duration)
			time.Sleep(duration)
			data_frame[5] = byte(value)
			_, err := conn.Write(data_frame)
			if err != nil {
				Debug("PrefixGenerator: cannot write byte", "value", value)
			}
			// here in future must read the response of coordinator to trottle or slowdonw
			Debug("PrefixGenerator: emitted payload", "duration", duration, "value", value)
		}
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 64)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			Debug("Sink: read  bytes from the connection ", "n", n, "byte", buf[0])

		} else if err != nil {
			Debug(" Sink: Can't read from connection")
			return
		}
	}

}
func testSinkGeneratorLoopTCP(coordinatorAddress string, port int) func() {

	return func() {
		registrationData := []byte{
			0x53, 0x54, // "ST"
			0x00, 0x05, // payload length = 5
			0x01,       // REGISTER
			0x00, 0x01, // TopologyID = 1
			0x00, 0x04, // NodeID = 4
		}
		registered := []byte{
			0x53, 0x54, // "ST"
			0x00, 0x01, // payload length = 1
			0x02, // REGISTERED
		}
		coordinatorData := make([]byte, len(registered))

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", coordinatorAddress, port))

		if err != nil {
			Debug("Sink: Can't open the connection to Coordinator  ", "hostname", coordinatorAddress, "port", port)
			log.Fatal("Quitting ..")
		}
		defer conn.Close()

		wb, err := conn.Write(registrationData)
		if err != nil {
			Debug("Cannot Write to Coordinator ", "address", coordinatorAddress)
			return
		}
		Debug("Written to Coordinator ", "bytes", wb)

		rb, err := io.ReadFull(conn, coordinatorData)
		if err != nil {
			Debug("Cannot Read From Coordinator ", "error", err)
			return
		}
		Debug("Read From Coordinator ", "bytes", rb)
		if !bytes.Equal(coordinatorData, registered) {
			Debug("Sink: unexpected registration response", "data", coordinatorData)
			return
		}
		Debug("Sink: registered in Coordinator", "topology_id", 1, "node_id", 4)

		handleConnection(conn)
	}
}
