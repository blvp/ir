package main
import (
	"strings"
	"io/ioutil"
	"regexp"
	"fmt"
)
func main() {
	fmt.Println("Hello")
}
func ObtainFile(filePath string) ([]string, map[string]int, map[string][]int) {
	replacer := strings.NewReplacer(
		",", "",
		":", "",
		";", "",
		"\"", "",
		"'", "",
		"-", "",
		"[", "",
		"]", "",
		"“", "",
		"”", "",
	)
	file, _ := ioutil.ReadFile(filePath)
	entireFile := strings.ToLower(string(file))
	entireFile = replacer.Replace(entireFile)

	documents := regexp.MustCompile("[.!?]").Split(entireFile, -1)
	stemmer := regexp.MustCompile("[[:space:]]")
	uniqueWords := map[string]string{}
	wordFreq := map[string]int{}
	postingLists := map[string][]int{}

	for docId, doc := range documents {
		words := stemmer.Split(doc, -1)
		for _, word := range words {
			if len([]rune(word)) != 0 {
				uniqueWords[word] = word
				wordFreq[word] += 1
				postingList := postingLists[word]

				if postingList != nil {
					postingList = append(postingList, docId)
				} else {
					postingList = []int{docId}
				}
				postingLists[word] = postingList
			}
		}
	}

	words := make([]string, 0, len(uniqueWords))

	for _, word := range uniqueWords {
		words = append(words, word)
	}
	return words, wordFreq, postingLists
}
