package main
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func StubDictionary() *Dictionary {
	words := []string{
		"test",
		"myday",
		"word",
		"yo",
		"machalka",
		"change",
		"oook",
		"love",
	}
	return StubDictionaryWords(words...)
}

func StubDictionaryWords(words ...string) *Dictionary {
	wordFreq := map[string]int{}
	for _, word := range words {
		wordFreq[word] += 1
	}

	return NewDictionary(words, wordFreq, 4)
}

func TestSplitIntoChunks(t *testing.T) {
	assert.Equal(t, 3, len(splitIntoChunks([]string{
		"test",
		"myday",
		"word",
		"yo",
		"yo",
		"machalka",
		"machalka",
		"change",
		"oook",
		"love",
	}, 4)))
}

func TestDictionaryCreate(t *testing.T) {
	dictionary := StubDictionary()
	assert.Equal(t, "6*change4&love8&machalka5&myday4*oook4&test4&word2&yo", dictionary.DictAsString)
	assert.Equal(t, 0, dictionary.PtrBlock[0].Ptr)
	assert.Equal(t, 31, dictionary.PtrBlock[1].Ptr)
}

func TestDictionaryTermLookup(t *testing.T) {
	dictionary := StubDictionary()
	lookup := dictionary.TermLookup("test")
	assert.Equal(t, 1, lookup.DocFreq)
	assert.Empty(t, lookup.PostingListPtr)


	dict := StubDictionaryWords(
		"testa",
		"testba",
		"testbb",
		"testc",
		"aaaa",
		"bbbb",
		"cccc",
		"dddd",
	)

	assert.Equal(t, 1, dict.TermLookup("testba").DocFreq)
}

func TestDictionaryFindWord(t *testing.T) {
	assert.Equal(t, "*change", StubDictionary().FindBlockHeader(0))
	assert.Equal(t, "*oook", StubDictionary().FindBlockHeader(31))
}

func TestDictionaryDecodeBlock(t *testing.T) {
	dict := StubDictionaryWords(
		"testa",
		"testaa",
		"testba",
		"testc",
		"aaaa",
		"bbbb",
		"cccc",
		"dddd",
	)
	wordsInBlock := dict.DecodeBlock(dict.PtrBlock[1])
	assert.Equal(t, 4, len(wordsInBlock))
	assert.Equal(t, []string{"testa", "testaa", "testba", "testc"}, wordsInBlock)
}

func TestObtainFile(t *testing.T) {
	words, wordFreq := ObtainFile("test_test.txt")
	assert.NotEmpty(t, words)
	assert.NotEmpty(t, wordFreq)
	assert.Equal(t, 8, len(words))
	assert.Equal(t, 8, len(wordFreq))
}

func TestCompress(t *testing.T) {

	testCases := map[string][]string{
		//prefixed fully first word
		"1a*2&a2&b3&cc": []string{
			"a",
			"aa",
			"ab",
			"acc",
		},
		//withoutPrefix
		"4*aaaa4&bbbb4&cccc4&dddd": []string{
			"aaaa",
			"bbbb",
			"cccc",
			"dddd",
		},
		//prefixed part of first word
		"4aa*aa4&bb4&cc4&dd": []string{
			"aaaa",
			"aabb",
			"aacc",
			"aadd",
		},
	}
	for expected, words := range testCases {
		assert.Equal(t, expected, Compress(words))
	}

}

type StringPair struct {
	first, second string
}

func NewStringPair(first, second string) StringPair {
	return StringPair{first:first, second: second}
}

func TestMatchIndex(t *testing.T) {

	testCases := map[StringPair]int{
		NewStringPair("aa", "cc"): 0,
		NewStringPair("aa", "aab"): 2,
		NewStringPair("aaaa", "aa"): 2,
		NewStringPair("", ""): 0,
		NewStringPair("", "asd"): 0,
		NewStringPair("asd", ""): 0,
	}
	for pair, expected := range testCases {
		assert.Equal(t, expected, matchIndex(pair.first, pair.second))
	}
}

func TestForSmallCollection(t *testing.T) {
	words, wordFreq := ObtainFile("test.txt")
	dict := NewDictionary(words, wordFreq, 4)
	fmt.Println(dict.DictAsString)
//	assert.Equal(t, 2, dict.TermLookup("will").DocFreq)
	for k, v := range wordFreq {
		fmt.Println(k)
		word := dict.TermLookup(k)
		assert.NotEmpty(t, word)
		assert.Equal(t, v, word.DocFreq)
	}

}