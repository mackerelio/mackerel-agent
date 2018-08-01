package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type graph struct {
	Label   string   `json:"label"`
	Unit    string   `json:"unit"`
	Metrics []metric `json:"metrics"`
}

type metric struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

var graphDefinition = graph{
	Label: "My Dice",
	Unit:  "integer",
	Metrics: []metric{
		{
			Name:  "d6",
			Label: "Die (d6)",
		},
		{
			Name:  "d20",
			Label: "Die (d20)",
		},
	},
}

func main() {
	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") == "1" {
		meta := struct {
			Graphs map[string]graph `json:"graphs"`
		}{
			Graphs: map[string]graph{"dice": graphDefinition},
		}
		bs, _ := json.Marshal(meta)
		fmt.Println("# mackerel-agent-plugin")
		fmt.Printf("%s\n", string(bs))
		os.Exit(0)
	}
	fmt.Printf("%s\t%d\t%d\n", "dice.d6", rand.Int()%6+1, time.Now().Unix())
	fmt.Printf("%s\t%d\t%d\n", "dice.d20", rand.Int()%20+1, time.Now().Unix())
}
