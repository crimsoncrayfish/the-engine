package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    fmt.Print("Press enter to close...")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}
