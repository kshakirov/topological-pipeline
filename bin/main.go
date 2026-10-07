package main

func main() {
	Debug("Starting Go Storm")
	generator := PrefixGeneratorTCP{Func: testPrefixGeneratorLoopTCP("127.0.0.1", 9300)}
	consumer := SinkTCP{Func: testSinkGeneratorLoopTCP("127.0.0.1", 9300)}

	// parallelNode := ParallelComposition(TcpSplitNodeConfig{ TcpExternalChoiceConfig{"127.0.0.1",9090},
	// 	TcpSpliterBuffferConfig{"127.0.0.1",9091}, 3}, testComputeB)
	// parallelNode.Start()

	generator.Start()
	consumer.Consume()
	//Debug("Go Storm stopped gracefully")
}
