package main
import (
	"regexp"
	"github.com/bradfitz/iter"
	"unicode/utf8"
	"strconv"
	"sort"
	"math"
)

type Dictionary struct {
	PtrBlock     []Block
	BlockSize    int
	DictAsString string
}

type Block struct {
	Ptr   int
	Words []Word
}

type Word struct {
	DocFreq        int
	PostingListPtr *int // Id of documents or Inverse Document Freq
}

func (d Dictionary) TermLookup(term string) *Word {
	l := 0
	r := len(d.PtrBlock)
	foundBlockHeader := ""
	for l < r {
		middle := (r + l) / 2
		ptr := d.PtrBlock[middle].Ptr
		foundBlockHeader = d.FindBlockHeader(ptr)
		if foundBlockHeader < term {
			l = middle + 1
		} else {
			r = middle - 1
		}
	}

	block := d.PtrBlock[int(math.Min(float64(l), float64(len(d.PtrBlock) - 1)))]
	words := d.DecodeBlock(block)

	for i, word := range words {
		if word == term {
			return &block.Words[i]
		}
	}
	return nil
}

func (d Dictionary) DecodeBlock(block Block) []string {
	result := []string{}
	ptr := block.Ptr
	for _ = range iter.N(d.BlockSize) {
		word := d.FindBlockHeader(ptr)
		result = append(result, word)
		wordLen := utf8.RuneCountInString(word)
		wordLenStr := strconv.Itoa(wordLen)
		ptr += utf8.RuneCountInString(wordLenStr) + wordLen
	}

	return result
}

func (d Dictionary) FindBlockHeader(ptr int) string {
	rightSide := []rune(d.DictAsString)[ptr: utf8.RuneCountInString(d.DictAsString)]
	wordLenStr := regexp.MustCompile("[0-9]+").FindString(string(rightSide))
	wordLen, _ := strconv.Atoi(wordLenStr)
	wordLenStrLen := utf8.RuneCountInString(wordLenStr)
	return string(rightSide[wordLenStrLen: wordLen + wordLenStrLen])
}


func NewDictionary(words []string, docFreq map[string]int, blockSize int) *Dictionary {
	sort.Strings(words)
	buffer := ""
	blocks := []Block{}

	for _, yo := range splitIntoChunks(words, blockSize) {
		wordInBlock := make([]Word, 0, len(yo))
		blockPtr := utf8.RuneCountInString(buffer)
		for _, word := range yo {
			wordInBlock = append(wordInBlock, Word{DocFreq:docFreq[word], PostingListPtr: nil})
			buffer += strconv.Itoa(utf8.RuneCountInString(word))
			buffer += word
		}
		blocks = append(blocks, Block{Ptr: blockPtr, Words: wordInBlock})
	}

	return &Dictionary{DictAsString: buffer, BlockSize:blockSize, PtrBlock:blocks}

}

func splitIntoChunks(arr []string, chunkSize int) [][]string {
	arrLen := len(arr)
	chunkNum := arrLen / chunkSize + 1
	result := make([][]string, 0, chunkNum)
	for i := 0; i < chunkNum; i += 1 {
		if i * chunkSize >= arrLen {
			break
		}
		tempSize := (i + 1) * chunkSize
		if tempSize >= arrLen {
			tempSize = arrLen
		}
		result = append(result, []string(arr[i * chunkSize: tempSize]))
	}
	return result;
}
