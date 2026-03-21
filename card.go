package main

import "fmt"

type Suit int

const (
	Spades Suit = iota
	Hearts
	Clubs
	Diamonds
)

type Rank int

const (
	Ace Rank = iota + 1
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

type Card struct {
	Rank Rank
	Suit Suit
}

// Value returns the cribbage counting value
func (c Card) Value() int {
	if c.Rank >= Ten {
		return 10
	}
	return int(c.Rank)
}

var suitNames = map[Suit]string{
	Spades:   "Spades",
	Hearts:   "Hearts",
	Clubs:    "Clubs",
	Diamonds: "Diamonds",
}

var rankNames = map[Rank]string{
	Ace:   "Ace",
	Two:   "Two",
	Three: "Three",
	Four:  "Four",
	Five:  "Five",
	Six:   "Six",
	Seven: "Seven",
	Eight: "Eight",
	Nine:  "Nine",
	Ten:   "Ten",
	Jack:  "Jack",
	Queen: "Queen",
	King:  "King",
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", rankNames[c.Rank], suitNames[c.Suit])
}
