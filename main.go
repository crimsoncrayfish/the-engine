package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("Press enter to close...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Print("test")
	if "test" == "test" {
		fmt.Print("Test")
	}
}
