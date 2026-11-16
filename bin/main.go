package main

func main() {
	Debug("Starting Go Storm")
	builder := Builder{Id: 1}
	topology := []byte{
		0x53, 0x54, // ST
		0x00, 0x15, // Payload length: 21 bytes
		0x03,       // CREATE_TOPOLOGY
		0x00, 0x01, // TopologyID: 1
		0x00, 0x03, // NodeCount: 3
		0x00, 0x01, // NodeID: 1
		0x00, 0x02, // NodeID: 2
		0x00, 0x03, // NodeID: 3
		0x00, 0x02, // EdgeCount: 2
		0x00, 0x01, // FromID: 1
		0x00, 0x02, // ToID: 2
		0x00, 0x02, // FromID: 2
		0x00, 0x03, // ToID: 3
	}
	builder.RegisterTopology("127.0.0.1:9300", topology)
	generator := PrefixGeneratorTCP{Func: testPrefixGeneratorLoopTCP("127.0.0.1", 9300)}
	consumer := SinkTCP{Func: testSinkGeneratorLoopTCP("127.0.0.1", 9300)}

	parallelNode := ParallelComposition(TcpSplitNodeConfig{TcpExternalChoiceConfig{"127.0.0.1", 9090},
		TcpSpliterBuffferConfig{"127.0.0.1", 9091}, 3}, testComputeB)
	parallelNode.Start()

	generator.Start()
	consumer.Consume()
	//Debug("Go Storm stopped gracefully")
}
