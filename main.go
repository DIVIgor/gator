package main

import (
	"fmt"
	"log"

	"github.com/DIVIgor/gator/internal/config"
)


func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("lane")
	

	cfg, err = config.Read()
	if err != nil {
		log.Fatal("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}