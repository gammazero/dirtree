package dirtree

import "sort"

// nodeSlice is used for sorting.
type nodeSlice []*Dirent

// Sort alphabetically sorts (from a to z) the slice of nodes.
func Sort(dSlice []*Dirent) {
	sort.Sort(nodeSlice(dSlice))
}

// Sort reverse-alphabetically sorts (from z to a) the slice of nodes.
func SortReverse(dSlice []*Dirent) {
	sort.Sort(sort.Reverse(nodeSlice(dSlice)))
}

func (s nodeSlice) Less(i, j int) bool { return s[i].String() < s[j].String() }
func (s nodeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nodeSlice) Len() int           { return len(s) }
