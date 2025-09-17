package main

import (
	"fmt"

	"github.com/vavarine/ttq/cmd"
)

var version string = "dev"

func main() {
	fmt.Printf("ttq version: %s\n", version)

	cmd.Execute()
}
