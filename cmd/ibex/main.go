package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	path = ""
	port = "31415"
)

func init() {
	fmt.Println(" _____  _               ")
	fmt.Println(" \\_    \\ |__   _____  __")
	fmt.Println("   / /\\/ '_ \\ / _ \\ \\/ /")
	fmt.Println("/\\/ /_ | |_) |  __/>  < ")
	fmt.Println("\\____/ |_.__/ \\___/_/\\_\\")
	fmt.Println("           v0.1")
	flag.StringVar(&path, "path", path, "Path to share")
	flag.StringVar(&port, "port", port, "Port to bind")
}

func main() {
	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		os.Exit(127)
		//wd, err := os.Getwd()
		//if err != nil {
		//	fmt.Printf("Error getting current path %v\n", err)
		//	os.Exit(127)
		//}
		//path = wd
	}
	fmt.Printf("[ibex]: Sharing this: %s\n", path)
	server := NewServer(":"+port, path)
	server.Start()
}
