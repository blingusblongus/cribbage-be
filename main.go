package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
)

func main() {
	http.HandleFunc("/api/handStats", withMetrics("handStats", handleHandStats))
	// http.HandleFunc("/api/random", handleRandom)
	http.HandleFunc("/metrics", metricsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}

func handleHandStats(w http.ResponseWriter, r *http.Request) {
	cardsParam := r.URL.Query().Get("cards")
	if cardsParam == "" {
		writeError(w, 400, "No cards provided")
		return
	}

	cards, err := ParseCards(cardsParam)
	if err != nil {
		writeError(w, 400, "Error parsing cards: "+err.Error())
		return
	}

	if len(cards) < 4 || len(cards) > 6 {
		writeError(w, 400, "Provide 4 to 6 cards")
		return
	}

	detail := strings.Split(r.URL.Query().Get("detail"), ",")
	sortBy := r.URL.Query().Get("sort")

	showAll := slices.Contains(detail, "all")
	showFlips := slices.Contains(detail, "flips")

	if len(cards) == 4 {
		result := AnalyzeHand(cards, nil)
		writeJSON(w, formatAnalysis(result, showFlips))
		return
	}

	// 5-6 cards: evaluate keep/discard combos
	type combo struct {
		Keep    []string        `json:"keep"`
		Discard []string        `json:"discard"`
		Result  json.RawMessage `json:"result"`
		mean    float64
		max     int
		min     int
	}

	var results []combo
	for _, keep := range combinations(cards, 4) {
		discard := difference(cards, keep)
		analysis := AnalyzeHand(keep, discard)

		formatted, _ := json.Marshal(formatAnalysis(analysis, showFlips))

		results = append(results, combo{
			Keep:    cardNames(keep),
			Discard: cardNames(discard),
			Result:  formatted,
			mean:    analysis.Mean,
			max:     analysis.Max,
			min:     analysis.Min,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		switch sortBy {
		case "max":
			return results[i].max > results[j].max
		case "min":
			return results[i].min > results[j].min
		default:
			return results[i].mean > results[j].mean
		}
	})

	if !showAll {
		results = results[:1]
	}

	writeJSON(w, results)
}

func handleRandom(w http.ResponseWriter, r *http.Request) {
	// TODO: Not impl'd
}

func formatAnalysis(a AnalysisResult, showFlips bool) map[string]any {
	options := make(map[string]any)

	for score, opt := range a.ScoringOptions {
		entry := map[string]any{
			"count":  opt.Count,
			"chance": opt.Chance,
		}

		if showFlips {
			entry["flips"] = opt.Flips
		}

		options[fmt.Sprintf("%d", score)] = entry
	}

	return map[string]any{
		"mean":              a.Mean,
		"max":               a.Max,
		"min":               a.Min,
		"standardDeviation": a.StandardDeviation,
		"scoringOptions":    options,
	}
}

func difference(all, exclude []Card) []Card {
	ex := make(map[Card]bool)

	for _, c := range exclude {
		ex[c] = true
	}

	var result []Card
	for _, c := range all {
		if !ex[c] {
			result = append(result, c)
		}
	}

	return result
}

func cardNames(cards []Card) []string {
	names := make([]string, len(cards))
	for i, c := range cards {
		names[i] = c.String()
	}
	return names
}
