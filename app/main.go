package main

import "fmt"
import "net"


const configPath = "."
const configFile = "properties"
const configExtension = "env" 

func main() {
	config.LoadConfig(configPath, configFile, configExtension)
	
	fmt.Println("Hello")
}

func connect() {
	tcpAddr := "wss://ws.xtb.com/demo"
	username := "adam.kochanski.97@gmail.com"
	password := 
	net.ResolveTCPAddr()
}
