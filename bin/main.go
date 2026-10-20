package main

func main() {
	Debug("Starting Go Storm")
	builder := Builder{Id: 1}
	topology := []byte{0x53, 0x54, 0x00, 0x03, 0x03, 0x00, 0x01}
	builder.RegisterTopology("127.0.0.1:9300", topology)
	generator := PrefixGeneratorTCP{Func: testPrefixGeneratorLoopTCP("127.0.0.1", 9300)}
	consumer := SinkTCP{Func: testSinkGeneratorLoopTCP("127.0.0.1", 9300)}

	// parallelNode := ParallelComposition(TcpSplitNodeConfig{ TcpExternalChoiceConfig{"127.0.0.1",9090},
	// 	TcpSpliterBuffferConfig{"127.0.0.1",9091}, 3}, testComputeB)
	// parallelNode.Start()

	generator.Start()
	consumer.Consume()
	//Debug("Go Storm stopped gracefully")
}
