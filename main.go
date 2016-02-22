package main
import (
	"fmt"
	"strings"
	"io/ioutil"
	"regexp"
)


func main() {
	words, wordFreq := ObtainFile("test.txt")
	dict := NewDictionary(words, wordFreq, 4)
	fmt.Println(dict.TermLookup("сказал"))
}

func ObtainFile(filePath string) ([]string, map[string]int) {
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
	for _, doc := range documents {
		words := stemmer.Split(doc, -1)
		for _, word := range words {
			if len([]rune(word)) != 0 {
				uniqueWords[word] = word
				wordFreq[word] += 1
			}
		}
	}
	words := make([]string, 0, len(uniqueWords))

	for _, word := range uniqueWords {
		words = append(words, word)
	}
	return words, wordFreq
}
