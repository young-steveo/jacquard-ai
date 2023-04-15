package tools

import (
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
)

type Sentiment int

const (
	Positive Sentiment = iota
	Negative
	Ambiguous
)

var positives = []string{"y", "yes", "ye"}
var negatives = []string{"n", "no", "nope"}

func GetSentiment(text string) (Sentiment, error) {
	// shortcut for small words
	// if text matches a positive or negative word, return that sentiment
	for _, positive := range positives {
		if text == positive {
			return Positive, nil
		}
	}
	for _, negative := range negatives {
		if text == negative {
			return Negative, nil
		}
	}

	// otherwise, use vader-go to parse the text and return the sentiment

	parsed := sentitext.Parse(text, lexicon.DefaultLexicon)
	sentiment := sentitext.PolarityScore(parsed)

	if sentiment.Compound > 0.05 {
		return Positive, nil
	} else if sentiment.Compound < -0.05 {
		return Negative, nil
	}
	return Ambiguous, nil

}
