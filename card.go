//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit represent the suit of the card
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Heart
	Club
	Joker
)

var suits = [...]Suit{Spade, Diamond, Heart, Club}

// Rank represent the rank of the card
type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
}

// New create new deck of card with options
func New(opts ...func([]Card) []Card) []Card {
	var card []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			card = append(card, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		card = opt(card)
	}
	return card
}

// DefaultSort sort the deck ascendingly
func DefaultSort(c []Card) []Card {
	sort.Slice(c, less(c))
	return c
}

// Sort custom sort the deck by the provided l function
func Sort(l func(c []Card) func(i, j int) bool) func(c []Card) []Card {
	return func(c []Card) []Card {
		sort.Slice(c, l(c))
		return c
	}
}

func less(c []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(c[i]) < absRank(c[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// Shuffle suffle a deck of card and return the new shuffled deck
func Shuffle(c []Card) []Card {
	ret := make([]Card, len(c))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(c))
	for i, j := range perm {
		ret[i] = c[j]
	}
	return ret
}

// Jokers create a deck with n Joker cards
func Jokers(n int) func([]Card) []Card {
	return func(c []Card) []Card {
		for i := 0; i < n; i++ {
			c = append(c, Card{Suit: Joker, Rank: Rank(i)})
		}
		return c
	}
}

// Filter filter the deck by function f
func Filter(f func(c Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, card := range cards {
			if !f(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

// Deck create a function that return a multiple decks
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
