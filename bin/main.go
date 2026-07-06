package main

import (
	"fmt"
	sp "gostorm.org/go_storm/lib/spout"
	bt "gostorm.org/go_storm/lib/bolt"
	tpl "gostorm.org/go_storm/lib/tuple"
)

func main() {
	fmt.Printf("")
	var spout sp.Spout
	var bolt bt.Bolt
	var touple tpl.Tuple 
	spout.Id = 1
	spout.NextTuple()
	bolt.Id = 2
	touple.Id = "3"
	touple.Add("Hello")
	fmt.Printf("%v\n", touple)
	bolt.Execute(&touple)
}
