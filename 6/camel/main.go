package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string
	fmt.Scanf("%s\n", &input)
	answer := 1
	for _, ch := range input {
		// min, max := 'A', 'Z'
		// if ch >= min && ch <= max {
		// 	answer++
		// }
		str := string(ch)
		if strings.ToUpper(str) == str {
			answer++
		}
	}

	fmt.Println(answer)
}
