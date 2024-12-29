package src

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
