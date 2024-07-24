package main

import (
	"fmt"
	"log"
)

func main() {
	// serve()
	
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	
	for i, host := range config.Hosts {
		fmt.Printf("Host %d:\n", i+1)
		fmt.Println("  Host:", host.Host)
		fmt.Println("  User:", host.User)
		if host.Port != nil {
			fmt.Println("  Port:", *host.Port)
		} else {
			fmt.Println("  Port is nil")
		}
		if host.Password != nil {
			fmt.Println("  Password:", *host.Password)
		} else {
			fmt.Println("  Password is nil")
		}
		if host.Identity != nil {
			fmt.Println("  Identity:", *host.Identity)
		} else {
			fmt.Println("  Identity is nil")
		}
	}
}
