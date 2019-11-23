package main

import (
	"fmt"
	"strings"

	"./deck"
)

type Hand []deck.Card

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
		State:  gs.State,
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

// CurrentPlayer return the hand of the current player
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("")
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// DealerString print dealer card, hide the second card
func (h Hand) DealerString() string {
	return h[0].String() + ", HIDDEN"
}

// MinScore return the min score of a hand
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Score return the score of a hand
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// BlackJack check if a hand is black jack
func BlackJack(hand Hand) bool {
	var hasAce, hasTen = false, false
	for _, h := range hand {
		if h.Rank == deck.Ace {
			hasAce = true
		} else if int(h.Rank) >= 10 {
			hasTen = true
		} else {
			return false
		}
	}
	return hasAce && hasTen
}

func shuffle(gs GameState) GameState {
	res := clone(gs)
	res.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return res
}

func deal(gs GameState) GameState {
	res := clone(gs)
	res.Player = make(Hand, 0, 5)
	res.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, res.Deck = draw(res.Deck)
		res.Player = append(res.Player, card)
		card, res.Deck = draw(res.Deck)
		res.Dealer = append(res.Dealer, card)
	}
	res.State = StatePlayerTurn
	return res
}

func hit(gs GameState) GameState {
	res := clone(gs)
	hand := res.CurrentPlayer()
	var card deck.Card
	card, res.Deck = draw(res.Deck)
	*hand = append(*hand, card)
	if hand.Score() >= 25 {
		return stand(res)
	}
	return res
}

func stand(gs GameState) GameState {
	res := clone(gs)
	res.State++
	return res
}

func main() {
	var gs GameState
	gs = shuffle(gs)
	gs = deal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player: ", gs.Player)
		fmt.Println("Dealer: ", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it or (s)top")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			gs = hit(gs)
		case "s":
			fmt.Println("You decided to stop")
			gs = stand(gs)
		default:
			fmt.Scanf("Invalid option")
		}
	}

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = hit(gs)
		} else {
			gs = stand(gs)
		}
	}

	gs = annouce(gs)
}

func annouce(gs GameState) GameState {
	res := clone(gs)
	res.State++
	pScore, dScore := gs.Player.Score(), gs.Dealer.Score()
	fmt.Println("===Final Hands===")
	fmt.Println("Player:", res.Player)
	fmt.Println("Player Score:", pScore)
	fmt.Println("Dealer:", res.Dealer)
	fmt.Println("Dealer Score:", dScore)

	switch {
	case pScore > 21:
		fmt.Println("Player busted")
	case dScore > 21:
		fmt.Println("Dealer Busted")
	default:
		if pScore > dScore {
			fmt.Println("Player wins")
		} else if dScore > pScore {
			fmt.Println("Dealer wins")
		} else {
			fmt.Println("Draw")
		}
	}
	res.Player = nil
	res.Dealer = nil
	return res
}
