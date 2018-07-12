package main

import (
	"fmt"

	"chirp.com/utils"

	"chirp.com/pkg/unique"
)

func main() {
	s := []string{"hello", "a", "A", "HELLO", "hello"}
	fmt.Println(s)
	fmt.Println(unique.Strings(s, utils.NormalizeText))
}
