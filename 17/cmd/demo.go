package main

import (
	"fmt"
	secret "gophercises/17"
)

func main() {
	v := secret.Memory("my-test-key")
	err := v.Set("demo_key", "some random value")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain:", plain)
}
