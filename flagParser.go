package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var timeout time.Duration
var portList []int

var (
	host            = flag.String("host", "", "Hostname / IP of the target system")
	timeoutDuration = flag.Float64("timeout", 0.5, "Timeout in seconds")
	maxconn         = flag.Int("connections", 10, "Limit the maximum number of parallel connections ")
	ports           = flag.String("ports", "", "A port OR range of ports to scan Eg. 80 or 8000-8080 or 80,21,22")
)

func ParseFlags() {
	flag.Usage = func() {
		fmt.Println("Portunus -  A Minimalist Multithreaded Port Scanner")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	//No Arguments Provided
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Error: No Arguments Provided\n")
		flag.Usage()
	}

	//Hostname and port not provided
	if *host == "" || *ports == "" {
		flag.Usage()
	}
	// NOTE: Supports only one form of parsing at a time
	if strings.Contains(*ports, ",") && strings.Contains(*ports, "-") {
		fmt.Fprintln(os.Stderr, "Please Use comma seperated ports OR range of ports, not both at once")
		flag.Usage()
	}

	if strings.Contains(*ports, ",") {
		var portStrings = strings.Split(*ports, ",")
		for _, portNumber := range portStrings {
			port, err := strconv.Atoi(portNumber)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid Port Number "+portNumber, err)
				flag.Usage()
			}
			portList = append(portList, port)
		}
	}

	if strings.Contains(*ports, "-") {
		portList = append(portList, portRangeParser(*ports)...)
	} else {
		// Single port
		port, err := strconv.Atoi(*ports)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid Port Number "+*ports, err)
		}
		portList = append(portList, port)
	}

	var ms = int(*timeoutDuration * 1000)
	timeout = time.Duration(ms) * time.Millisecond
}

func portRangeParser(rangeString string) []int {
	portRange := strings.Split(rangeString, "-")
	var ports []int
	start, err := strconv.Atoi(portRange[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid Port at start of range", err)
		flag.Usage()
	}
	end, err := strconv.Atoi(portRange[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid Port at start of range", err)
	}

	for i := 0; i < end-start+1; i++ {
		ports = append(ports, start+i)
	}

	return ports
}
