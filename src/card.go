package src

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
)

type Card struct {
	FaceValue int64
	Picture   bool
}

type Player struct {
	Hand
	Bet   int64
	Money int64
	Buyin int64
	Id    int
}

func (p *Player) PercentBuyInChange() float64 {
	return (float64(p.Money-p.Buyin) / float64(p.Buyin)) * 100
}

func (p *Player) Buy(money int64) {
	p.Money += money
	p.Buyin += money
}

func (p *Player) Receive(money int64) {
	p.Money += money
}

func (p *Player) Give(money int64) {
	p.Receive(-money)
}

type Hand struct {
	Cards []*Card
}

type Deck = Hand

func (d *Hand) ShowCards(n int64) {
	sb := strings.Builder{}

	n = int64(min(int(n), len(d.Cards)))
	sb.WriteString(fmt.Sprintf("Hand til %d", n))

	for i := range n {
		sb.WriteString(fmt.Sprintf("%#v ", d.Cards[i]))
	}

	log.Printf(sb.String())
}

func (d *Hand) GetCards(n int64) []*Card {
	return d.Cards[0:n]
}

type Hander interface {
	Take(cards ...*Card)
	Discard()
	Shower
	CardGetter
}

type CardGetter interface {
	GetCards(n int64) []*Card
}

type Shower interface {
	ShowCards(n int64)
}

func (d *Hand) Take(cards ...*Card) {
	d.Cards = append(d.Cards, cards...)
}
func (d *Hand) Discard() {
	d.Cards = nil
}

func (d *Hand) Shuffle(n int) {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func NewDeck() Deck {
	var baseCards []Card
	for i := range 10 {
		i += 1
		baseCards = append(baseCards, Card{FaceValue: int64(i), Picture: false})
	}

	baseCards = append(baseCards, Card{FaceValue: 0, Picture: true}, Card{FaceValue: 0, Picture: true}, Card{FaceValue: 0, Picture: true})

	var cards []Card
	for range 4 {
		cards = append(cards, baseCards...)
	}

	pointyCards := make([]*Card, 0, len(cards))
	for i := range cards {
		pointyCards = append(pointyCards, &cards[i])
	}
	return Deck{Cards: pointyCards}
}
