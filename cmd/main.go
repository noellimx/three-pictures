package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"three-pictures/src"
)

const MAX_PLAYERS = 7
const DEFAULT_ROUNDS = 100000
const LOG_ROUND_EVERY = 10000

type Stats struct {
	DealerPctChange    float64  `json:"dealer_pct_change"`
	AvgPlayerPctChange float64  `json:"average_player_pct_change"`
	RoundsPlayed       int64    `json:"rounds_played"`
	Players            int      `json:"players"`
	HalfPayOn          src.RANK `json:"half_pay_on"`
}

var StatsList []Stats

func main() {
	sanityChecks()

	var wg sync.WaitGroup
	for rank := src.RANK_NONE; rank >= src.RANK_PICTURE; rank-- {
		for playerCount := 1; playerCount <= MAX_PLAYERS; playerCount++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				simulate(playerCount, rank, DEFAULT_ROUNDS)
			}()
		}
	}

	wg.Wait()

	csvData := strings.Builder{}
	csvData.WriteString("dealer_pct_change")
	csvData.WriteString(",average_player_pct_change")
	csvData.WriteString(",rounds_played")
	csvData.WriteString(",players")
	csvData.WriteString(",half_pay_on\n")
	for _, stats := range StatsList {
		csvData.WriteString(fmt.Sprintf("%f,%f,%d,%d,%d\n", stats.DealerPctChange, stats.AvgPlayerPctChange, stats.RoundsPlayed, stats.Players, stats.HalfPayOn))
	}

	f, err := os.OpenFile("some.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(csvData.String())
}

func simulate(playersCount int, halfPayOn src.RANK, roundsToPlay int64) {
	fmt.Sprintf("simulating....")
	var shouldHalfPay func(p *src.Player) bool
	if halfPayOn == src.RANK_PICTURE {
		shouldHalfPay = func(p *src.Player) bool {
			return src.IsTriplePicture(p)
		}
	} else {
		shouldHalfPay = func(p *src.Player) bool {
			return src.Points(p) == halfPayOn
		}
	}
	dealer := &src.Player{
		Id: 0,
	}
	dealer.Buy(10000)
	var seats []*src.Player
	//for p := range rand.IntN(MAX_PLAYERS) + 1 {
	for p := range playersCount {
		s := &src.Player{
			Id:   p + 1,
			Hand: src.Hand{},
			Bet:  10,
		}
		s.Buy(10000)
		seats = append(seats, s)
	}

	var dealerWins, playerWins, draws atomic.Int64

	var round atomic.Int64
	sb := strings.Builder{}
	logg := func() {
		defer sb.Reset()
		_de, _p, _dr := dealerWins.Load(), playerWins.Load(), draws.Load()
		totalHandsDealt := _de + _p + _dr

		percent := func(part, whole int64) string {
			return fmt.Sprintf("%f%%", float64(part)/float64(whole)*100)
		}
		_de_p, _p_p, _dr_p := percent(_de, totalHandsDealt), percent(_p, totalHandsDealt), percent(_dr, totalHandsDealt)

		_round := round.Load()

		var avgPlayersChange, totalPlayersChange float64
		for _, seat := range seats {
			totalPlayersChange += seat.PercentBuyInChange()
		}
		avgPlayersChange = totalPlayersChange / float64(len(seats))

		StatsList = append(StatsList, Stats{
			DealerPctChange:    dealer.PercentBuyInChange(),
			AvgPlayerPctChange: avgPlayersChange,
			RoundsPlayed:       _round,
			Players:            playersCount,
			HalfPayOn:          halfPayOn,
		})

		sb.WriteString(fmt.Sprintf("Round#%d -> dealer wins %d %s, player wins %d %s, draws %d %s, total hands dealt %d hands/round %2f \n", _round, _de, _de_p, _p, _p_p, _dr, _dr_p, totalHandsDealt, float64(totalHandsDealt)/float64(_round)))
		sb.WriteString(fmt.Sprintf("dealer: %8d$", dealer.Money))

		for i, s := range seats {
			sb.WriteString(fmt.Sprintf(", seat %d: %8d$", i, s.Money))
		}

		fmt.Println(sb.String())
	}

	for r := range roundsToPlay {
		round.Store(r + 1)

		Deal(src.NewDeck(), dealer, seats)

		for _, seat := range seats {
			winner, _ := src.CheckUpperHand(&dealer.Hand, &seat.Hand)
			if winner == nil {
				draws.Add(1)
				//log.Printf("dealer %#v seat %#v DRAW \n", src.Points(dealer), src.Points(seat))
			} else {
				var from, to *src.Player
				var bet int64
				if winner == &seat.Hand {
					playerWins.Add(1)

					bet = seat.Bet
					if shouldHalfPay(seat) {
						bet /= 2
					}
					from, to = dealer, seat
				} else {
					dealerWins.Add(1)

					bet = seat.Bet
					from, to = seat, dealer

				}
				Transfer(from, to, bet)
			}
		}

		ResetTable(dealer, seats)
		if r != 0 && (r+1)%LOG_ROUND_EVERY == 0 {
			logg()
		}
	}
}

func Transfer(from *src.Player, to *src.Player, bet int64) {
	from.Give(bet)
	to.Receive(bet)
}

// checkDeck
// check 52 cards, with total points
func checkDeck() {
	deck := src.NewDeck()

	// wantTotalPoints
	// all cards with face value
	var wantTotalPoints int64 = (1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10) * 4

	var gotTotalPoints int64 = 0
	for _, c := range deck.Cards {
		gotTotalPoints += c.FaceValue
	}

	if wantTotalPoints != gotTotalPoints {
		log.Panicf("points mismatch %d %d \n", wantTotalPoints, gotTotalPoints)
	}

	if len(deck.Cards) != 52 {
		log.Panicf("length mismatch \n")
	}
}

func sanityChecks() {
	checkDeck()
	checkRank()
}

func checkRank() {
	triples := &src.Hand{
		Cards: []*src.Card{
			{
				FaceValue: 0,
				Picture:   true,
			}, {
				FaceValue: 0,
				Picture:   true,
			}, {
				FaceValue: 0,
				Picture:   true,
			},
		},
	}

	for i := range 10 {
		for j := range 10 {
			i += 1
			j += 1
			winner, loser := src.CheckUpperHand(triples, &src.Hand{
				Cards: []*src.Card{
					{
						FaceValue: int64(i),
						Picture:   false,
					}, {
						FaceValue: int64(j),
						Picture:   false,
					}, {
						FaceValue: 0,
						Picture:   true,
					},
				},
			})

			if winner != triples {
				log.Panicf("winner should be triples %#v %#v\n", winner, triples)
			}

			if int(src.Points(loser)) != (i+j)%10 {
				log.Panicf("loser should have %d points, got%#v \n", (i+j)%10, src.Points(loser))
			}
		}
	}
}
