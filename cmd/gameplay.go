package main

import (
	"math/rand/v2"

	"three-pictures/src"
)

// ResetTable discards all cards
func ResetTable(dealer src.Hander, players []*src.Player) {
	dealer.Discard()
	for _, p := range players {
		p.Discard()
	}
}

func Deal(deck src.Deck, dealer src.Hander, players []*src.Player) {
	deck.Shuffle(52)
	var hands []src.Hander
	hands = append(hands, dealer)
	for _, player := range players {
		hands = append(hands, player)
	}
	rand.Shuffle(len(hands), func(i, j int) {
		hands[i], hands[j] = hands[j], hands[i]
	})
	for i, h := range hands {
		for j := range 3 {
			h.Take(deck.Cards[i*2+j])
		}
	}
}
