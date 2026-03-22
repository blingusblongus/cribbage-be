package main

type ScoringDetail struct {
	Category string
	Cards    []Card
	Points   int
}

type ScoreResult struct {
	Details []ScoringDetail
}

func (r ScoreResult) Total(includeNibs bool) int {
	total := 0
	for _, d := range r.Details {
		if !includeNibs && d.Category == "nibs" {
			continue
		}

		total += d.Points
	}

	return total
}

func ScoreHand(hand []Card, flip Card) ScoreResult {
	all := make([]Card, len(hand)+1)
	copy(all, hand)
	all[len(hand)] = flip

	var details []ScoringDetail
	details = append(details, scoreFifteens(all)...)
	details = append(details, scorePairs(all)...)
	details = append(details, scoreRuns(all)...)
	details = append(details, scoreFlush(hand, flip)...)
	details = append(details, scoreNobs(hand, flip)...)
	details = append(details, scoreNibs(flip)...)

	return ScoreResult{Details: details}
}

func scoreFifteens(cards []Card) []ScoringDetail {
	var details []ScoringDetail
	for size := 2; size <= len(cards); size++ {
		for _, combo := range combinations(cards, size) {
			sum := 0
			for _, c := range combo {
				sum += c.Value()
			}
			if sum == 15 {
				details = append(details, ScoringDetail{
					Category: "fifteen",
					Cards:    combo,
					Points:   2,
				})
			}
		}
	}

	return details
}

func scorePairs(cards []Card) []ScoringDetail {
	var details []ScoringDetail
	for i := range cards {
		for j := i + 1; j < len(cards); j++ {
			if cards[i].Rank == cards[j].Rank {
				details = append(details, ScoringDetail{
					Category: "pair",
					Cards:    []Card{cards[i], cards[j]},
					Points:   2,
				})
			}
		}
	}

	return details
}

func scoreRuns(cards []Card) []ScoringDetail {
	for length := 5; length >= 3; length-- {
		var details []ScoringDetail
		for _, combo := range combinations(cards, length) {
			if isRun(combo) {
				details = append(details, ScoringDetail{
					Category: "run",
					Cards:    combo,
					Points:   length,
				})
			}
		}

		if len(details) > 0 {
			return details
		}
	}

	return nil
}

func isRun(cards []Card) bool {
	ranks := make([]int, len(cards))
	for i, c := range cards {
		ranks[i] = int(c.Rank)
	}

	for i := 1; i < len(ranks); i++ {
		for j := i; j > 0 && ranks[j-1] > ranks[j]; j-- {
			ranks[j], ranks[j-1] = ranks[j-1], ranks[j]
		}
	}

	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]+1 {
			return false
		}
	}

	return true
}

func scoreFlush(hand []Card, flip Card) []ScoringDetail {
	suit := hand[0].Suit

	for _, c := range hand[1:] {
		if c.Suit != suit {
			return nil
		}
	}

	if flip.Suit == suit {
		cards := make([]Card, len(hand)+1)
		copy(cards, hand)
		cards[len(hand)] = flip
		return []ScoringDetail{{
			Category: "flush",
			Cards:    cards,
			Points:   5,
		}}
	}

	return []ScoringDetail{{
		Category: "flush",
		Cards:    hand,
		Points:   4,
	}}
}

func scoreNobs(hand []Card, flip Card) []ScoringDetail {
	for _, c := range hand {
		if c.Rank == Jack && c.Suit == flip.Suit {
			return []ScoringDetail{{
				Category: "nobs",
				Cards:    []Card{c},
				Points:   1,
			}}
		}
	}

	return nil
}

func scoreNibs(flip Card) []ScoringDetail {
	if flip.Rank == Jack {
		return []ScoringDetail{{
			Category: "nibs",
			Cards:    []Card{flip},
			Points:   2,
		}}
	}

	return nil
}

func combinations(cards []Card, k int) [][]Card {
	var result [][]Card
	var build func(start int, current []Card)
	build = func(start int, current []Card) {
		if len(current) == k {
			combo := make([]Card, k)
			copy(combo, current)
			result = append(result, combo)
			return
		}

		for i := start; i < len(cards); i++ {
			build(i+1, append(current, cards[i]))
		}
	}
	build(0, nil)
	return result
}
