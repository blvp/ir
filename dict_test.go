package main
import (
	"testing"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "6change4love8machalka5myday4oook4test4word2yo", dictionary.DictAsString)
	assert.Equal(t, 0, dictionary.PtrBlock[0].Ptr)
	assert.Equal(t, 27, dictionary.PtrBlock[1].Ptr)
}

func TestDictionaryTermLookup(t *testing.T) {
	dictionary := StubDictionary()
	lookup := dictionary.TermLookup("test")
	assert.Equal(t, 1, lookup.DocFreq)
	assert.Empty(t, lookup.PostingListPtr)
}

func TestDictionaryFindWord(t *testing.T) {
	assert.Equal(t, "change", StubDictionary().FindBlockHeader(0))
	assert.Equal(t, "oook", StubDictionary().FindBlockHeader(27))
}

func TestDictionaryDecodeBlock(t *testing.T) {
	dictionary := StubDictionary()
	wordsInBlock := dictionary.DecodeBlock(dictionary.PtrBlock[0])
	assert.Equal(t, 4, len(wordsInBlock))
	assert.Equal(t, []string{"change", "love", "machalka", "myday"}, wordsInBlock)
}

func TestObtainFile(t *testing.T) {
	words, wordFreq := ObtainFile("test_test.txt")
	assert.NotEmpty(t, words)
	assert.NotEmpty(t, wordFreq)
	assert.Equal(t, 8, len(words))
	assert.Equal(t, 8, len(wordFreq))
}