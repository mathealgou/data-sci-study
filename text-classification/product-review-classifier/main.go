package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strings"
)

type Table []map[string]string

const THRESHOLD = 2

func readCsv(filePath string, columns []int) (Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	csv_reader := csv.NewReader(file)

	data, err := csv_reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		panic(err)

	}

	// return the result based on the columns indexes
	result := make(Table, len(data)-1)
	for i, row := range data[1:] {
		result[i] = make(map[string]string)
		for _, col := range columns {
			result[i][data[0][col]] = row[col]
		}
	}

	return result, nil

}

func generateBagOfWords(table Table) map[string]int {
	bow := make(map[string]int)

	for _, review := range table {
		reviewWords := tokenize(review["reviews.text"])
		for _, word := range reviewWords {
			val, ok := bow[word]
			if !ok {
				bow[word] = 1
			} else {
				bow[word] = val + 1
			}
		}
	}

	// applying a THRESHOLD
	result := make(map[string]int)
	for word, count := range bow {
		if count > THRESHOLD {
			result[word] = count
		}
	}

	return result
}

func tokenize(text string) []string {
	clean := strings.ToLower(text)
	clean = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		return ' '
	}, clean)
	return strings.Fields(clean)
}

func filterTable(table Table, column string, value string) Table {
	result := make(Table, 0)
	for _, row := range table {
		if row[column] == value {
			result = append(result, row)
		}
	}
	return result
}

func calculateTextBowProbability(text string, bow map[string]int, totalWords int) float64 {
	words := tokenize(text)
	probability := 0.0
	denominator := float64(totalWords + len(bow))

	for _, word := range words {
		count := bow[word]
		// Laplace smoothing for unseen words.
		probability += math.Log(float64(count+1) / denominator)
	}
	return probability
}

func calculateRecommendationProbability(testReview string, recommendBow, nonRecommendBow map[string]int, recommendCount, nonRecommendCount int) (float64, float64) {
	recommendPrior := float64(recommendCount) / float64(recommendCount+nonRecommendCount)
	nonRecommendPrior := float64(nonRecommendCount) / float64(recommendCount+nonRecommendCount)

	recommendTotal := 0
	for _, count := range recommendBow {
		recommendTotal += count
	}
	nonRecommendTotal := 0
	for _, count := range nonRecommendBow {
		nonRecommendTotal += count
	}

	logRecommend := math.Log(recommendPrior) + calculateTextBowProbability(testReview, recommendBow, recommendTotal)
	logNonRecommend := math.Log(nonRecommendPrior) + calculateTextBowProbability(testReview, nonRecommendBow, nonRecommendTotal)

	// Normalize via softmax to get comparable probabilities.
	maxLog := math.Max(logRecommend, logNonRecommend)
	expRecommend := math.Exp(logRecommend - maxLog)
	expNonRecommend := math.Exp(logNonRecommend - maxLog)

	norm := expRecommend + expNonRecommend
	return expRecommend / norm, expNonRecommend / norm
}

func main() {
	useful_columns := []int{0, 11, 16} // eu já falei que odeio essa sintaxe?

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide a review as an argument.")
		return
	}
	println("Review to analyze:", strings.Join(args, " "))

	fmt.Println("Parsing...")
	data, _ := readCsv("./data/part1.csv", useful_columns)

	fmt.Println("Training...")
	recommend := filterTable(data, "reviews.doRecommend", "TRUE")

	nonRecommend := filterTable(data, "reviews.doRecommend", "FALSE")

	recommendBow := generateBagOfWords(recommend)

	nonRecommendBow := generateBagOfWords(nonRecommend)

	fmt.Println("Testing...")
	testReview := strings.Join(args, " ")
	probRecommend, probNonRecommend := calculateRecommendationProbability(testReview, recommendBow, nonRecommendBow, len(recommend), len(nonRecommend))

	fmt.Printf("Probability of recommendation: %.4f\n", probRecommend)
	fmt.Printf("Probability of non-recommendation: %.4f\n", probNonRecommend)

	if probRecommend > probNonRecommend {
		fmt.Println("The review is likely to be a recommendation.")
	} else {
		fmt.Println("The review is likely to be a non-recommendation.")
	}
}
