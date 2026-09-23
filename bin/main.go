package main

func main() {
	Debug("Starting Go Storm")
	generator := PrefixGeneratorTCP{Func: testPrefixGeneratorLoopTCP("127.0.0.1", 9090)}
	consumer := SinkTCP{Func: testSinkGeneratorLoopTCP("127.0.0.1",9091)}

	
	ParallelComposition(TcpSplitNodeConfig{ TcpExternalChoiceConfig{"127.0.0.1",9090},
		TcpSpliterBuffferConfig{"127.0.0.1",9091}, 2}, testComputeB)
	

	generator.Start()
	consumer.Consume()
	Debug("Go Storm stopped gracefully")
}
