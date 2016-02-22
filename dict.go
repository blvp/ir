package main
import (
	"regexp"
	"github.com/bradfitz/iter"
	"unicode/utf8"
	"strconv"
	"sort"
	"math"
	"strings"
	"fmt"
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
	endOfPrefixReplacer := strings.NewReplacer("*", "")
	for l < r {
		middle := (r + l) / 2
		ptr := d.PtrBlock[middle].Ptr
		foundBlockHeader := d.FindBlockHeader(ptr)
		cleanHeader := endOfPrefixReplacer.Replace(foundBlockHeader)
		if cleanHeader < term {
			l = middle + 1
		} else if cleanHeader > term {
			r = middle - 1
		} else {
			r = middle
			l = middle
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
	prefix := d.ResolvePrefix(block)
	prefixReplacer := strings.NewReplacer("&", prefix, "*", "")
	prefixLen := utf8.RuneCountInString(prefix)
	for i := range iter.N(d.BlockSize) {
		//not always last block is block sized
		if ptr == utf8.RuneCountInString(d.DictAsString) {
			break
		}
		rightSide := []rune(d.DictAsString)[ptr:]
		wordLenStr := regexp.MustCompile("[0-9]+").FindString(string(rightSide))
		minusPrefix := prefixLen
		if i == 0 {
			minusPrefix = 0
		}
		wordLen, _ := strconv.Atoi(wordLenStr)
		wordLenStrLen := utf8.RuneCountInString(wordLenStr)
		result = append(result, prefixReplacer.Replace(
			// offset of wordLen and then we decide that & symbol is prefixLen weight
			string(rightSide[wordLenStrLen:wordLen + wordLenStrLen - minusPrefix + 1])))
		fmt.Println(result)
		ptr += wordLen + wordLenStrLen - minusPrefix + 1
	}

	return result
}add
func (d Dictionary) ResolvePrefix(block Block) string {
	ptr := block.Ptr
	rightSide := []rune(d.DictAsString)[ptr: ]
	wordLenStr := regexp.MustCompile("[0-9]+").FindString(string(rightSide))
	wordLenStrLen := utf8.RuneCountInString(wordLenStr)
	return string(rightSide[wordLenStrLen: strings.Index(string(rightSide), "*")])
}

func (d Dictionary) FindBlockHeader(ptr int) string {
	rightSide := []rune(d.DictAsString)[ptr: ]
	wordLenStr := regexp.MustCompile("[0-9]+").FindString(string(rightSide))
	wordLen, _ := strconv.Atoi(wordLenStr)
	wordLenStrLen := utf8.RuneCountInString(wordLenStr)
	word := string(rightSide[wordLenStrLen: wordLen + wordLenStrLen + 1])
	return word
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
		}
		buffer += Compress(yo)
		blocks = append(blocks, Block{Ptr: blockPtr, Words: wordInBlock})
	}

	return &Dictionary{DictAsString: buffer, BlockSize:blockSize, PtrBlock:blocks}

}

func Compress(words []string) string {
	firstWord := words[0]
	prefixLen := 1 << 32
	for _, word := range words {
		prefixLen = int(math.Min(float64(matchIndex(firstWord, word)), float64(prefixLen)))
	}
	prefix := string([]rune(firstWord)[0: prefixLen])
	firstWordLen := utf8.RuneCountInString(firstWord)
	buffer := strconv.Itoa(firstWordLen) + prefix + "*" + string([]rune(firstWord)[prefixLen:firstWordLen])
	for _, word := range words[1:] {
		wordLen := utf8.RuneCountInString(word)
		buffer += strconv.Itoa(wordLen)
		buffer += "&"
		buffer += string([]rune(word)[prefixLen: wordLen])
	}
	return buffer
}

func matchIndex(first, second string) int {
	//we should find an index of the second best match element
	//at least first length always lt second
	//so we should just iterate over and count until string does not match
	if first == "" || second == "" {
		return 0
	}
	firstArr := []rune(first)
	secondArr := []rune(second)
	n := int(math.Min(float64(len(firstArr)), float64(len(secondArr))))
	startPos := 0
	for i := 0; i < n; i += 1 {
		if string(firstArr[i]) == string(secondArr[i]) {
			startPos += 1
		} else {
			break;
		}
	}
	return startPos
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
