package main

import "math"

type FlipOutcome struct {
	Score  int
	Count  int
	Chance float64
	Flips  []string
}

type AnalysisResult struct {
	Mean              float64
	Max               int
	Min               int
	StandardDeviation float64
	ScoringOptions    map[int]*FlipOutcome
}

func AnalyzeHand(hand []Card, discarded []Card) AnalysisResult {
	exclude := make([]Card, 0, len(hand)+len(discarded))
	exclude = append(exclude, hand...)
	exclude = append(exclude, discarded...)
	remaining := RemainingCards(exclude)

	options := make(map[int]*FlipOutcome)
	total := 0
	max := 0
	min := math.MaxInt

	for _, flip := range remaining {
		score := ScoreHand(hand, flip).Total(false)
		total += score

		if score > max {
			max = score
		}
		if score < min {
			min = score
		}

		if opt, ok := options[score]; ok {
			opt.Count++
			opt.Flips = append(opt.Flips, flip.String())
		} else {
			options[score] = &FlipOutcome{
				Score: score,
				Count: 1,
				Flips: []string{flip.String()},
			}
		}
	}

	n := len(remaining)
	mean := float64(total) / float64(n)

	// stddev
	variance := 0.0
	for score, opt := range options {
		diff := float64(score) - mean
		variance += diff * diff * float64(opt.Count)
	}

	variance /= float64(n)
	stdDev := math.Sqrt(variance)

	// Chances
	for _, opt := range options {
		opt.Chance = float64(opt.Count) / float64(n) * 100
	}

	return AnalysisResult{
		Mean:              mean,
		Max:               max,
		Min:               min,
		StandardDeviation: stdDev,
		ScoringOptions:    options,
	}
}
