package main

import (
	"fmt"
	"strings"
)

var suitMap = map[byte]Suit{
	's': Spades,
	'h': Hearts,
	'c': Clubs,
	'd': Diamonds,
}

var rankMap = map[string]Rank{
	"a":  Ace,
	"2":  Two,
	"3":  Three,
	"4":  Four,
	"5":  Five,
	"6":  Six,
	"7":  Seven,
	"8":  Eight,
	"9":  Nine,
	"10": Ten,
	"j":  Jack,
	"q":  Queen,
	"k":  King,
}

func ParseCards(query string) ([]Card, error) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(query)), ",")
	seen := make(map[string]bool)
	cards := make([]Card, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) < 2 || len(p) > 3 {
			return nil, fmt.Errorf("invalid card: %q", p)
		}

		if seen[p] {
			return nil, fmt.Errorf("duplicate card: %q", p)
		}
		seen[p] = true

		suitChar := p[len(p)-1]
		rankStr := p[:len(p)-1]

		suit, ok := suitMap[suitChar]
		if !ok {
			return nil, fmt.Errorf("invalid suit in %q", p)
		}

		rank, ok := rankMap[rankStr]
		if !ok {
			return nil, fmt.Errorf("invalid rank in %q", p)
		}

		cards = append(cards, Card{Rank: rank, Suit: suit})
	}

	return cards, nil
}
