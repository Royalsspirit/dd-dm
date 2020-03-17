package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/Royalsspirit/dd-dm/internal/term"
)

func main() {
	file := flag.String("logfile", "./localfile.log", "a path of log file wanted to monitoring")
	threshold := flag.String("threshold", "10", "a threshold request per second")
	flag.Parse()

	fmt.Println("file", *file, "threshoold", *threshold)

	thresholdValue, _ := strconv.Atoi(*threshold)

	fmt.Println("intThres", thresholdValue)
	t := term.NewTerm(&term.Conf{
		Logfile:   *file,
		Threshold: thresholdValue,
	})

	t.Run()
}
