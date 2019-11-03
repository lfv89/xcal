package main

import (
	"flag"
	"fmt"
	"os"
	googleCalendar "xcal/internal/google"
)

var max int64 = 1
var truncate *string

func init() {
	nextCmd := flag.NewFlagSet("next", flag.ExitOnError)
	truncate = nextCmd.String("truncate", "20", "size of each event title")

	flag.Parse()
	err := nextCmd.Parse(os.Args[2:])

	if err != nil {
		panic("Error! Couldn't parse the args.")
	}
}

func main() {
	switch os.Args[1] {
	case "init":
		googleCalendar.Init()
	case "next":
		googleCalendar.GetNextEvent(max, truncate)
	default:
		fmt.Println("The command was not recognized.")
		os.Exit(1)
	}
}
