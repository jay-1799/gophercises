package main

import (
	"fmt"
	"gophercises/blackjack_ai/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
