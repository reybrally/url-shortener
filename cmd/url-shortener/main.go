package main

import (
	"fmt"
	"github.com/reybrally/REST-API-app/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)
}
