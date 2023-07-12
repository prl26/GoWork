package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "OmniFDN - Twitter Followers"
	fmt.Println(strings.Contains(s, "Twitter") && strings.Contains(s, "Follow"))
}
