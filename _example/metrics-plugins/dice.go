package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	prefix := flag.String("name", "example", "Metric name")
	flag.Parse()
	metricName, value := *prefix+".dice", rand.Int()%6+1
	fmt.Printf("%s\t%d\t%d\n", metricName, value, time.Now().Unix())
}
