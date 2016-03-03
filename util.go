package main
import "sort"

type Pair struct {
	a string
	b []int
}

type PairList []Pair
func (p PairList) Len() int { return len(p) }

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PairList) Less(i, j int) bool { return p[i].a < p[j].a }

func sortMap(a map[string][]int) PairList {

	result := make(PairList, 0, len(a))
	for k, v := range a {
		result = append(result, Pair{a: k, b:v})
	}
	sort.Sort(result)
	return result
}
