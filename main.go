package main

import (
	"fmt"
	"os"
)

func main() {
	path, _ := os.Getwd()
	fmt.Println(path)
}
