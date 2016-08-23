/*
Package dirtree provides basic traversal, search, and string conversion
functionality for a directory tree-like structure.

The propose of dirtree is to provide a way to construct a tree, that can be
traversed and printed, where each node is a simple container for its children.
This is useful for displaying and navigating containers whose structure can be
represented in a generalized way by this this package.

*/
package dirtree

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gammazero/queue"
)

// Dirent is an element in a tree similar to a basic directory tree.
type Dirent struct {
	name     string
	parent   *Dirent
	children map[string]*Dirent
}

// New creates a new root node.
func New(name string) *Dirent {
	return &Dirent{
		name: name,
	}
}

// String returns the name of the directory entry.
func (d *Dirent) String() string {
	if d.name == "" {
		return "/"
	}
	return d.name
}

// Size returns the number of children the receiver has.
func (d *Dirent) Size() int {
	return len(d.children)
}

// ForChild executes the given function for each child in the current Dirent.
func (d *Dirent) ForChild(f func(d *Dirent) bool) {
	for _, ch := range d.children {
		if !f(ch) {
			break
		}
	}
}

// ForParent executes the given function for each parent Dirent up to the root.
func (d *Dirent) ForParent(f func(d *Dirent) bool) {
	p := d.parent
	for p != nil {
		if !f(p) {
			break
		}
		p = p.parent
	}
}

// Child returns the Dirent with the given name, or nil if not found.
func (d *Dirent) Child(name string) *Dirent {
	ch, _ := d.children[name]
	return ch
}

// Add creates one or more new directory entry as a child of the receiver.
func (d *Dirent) Make(names ...string) error {
	var name string
	for _, name = range names {
		if name == "" || name == "/" {
			return fmt.Errorf("invalid name: %s", name)
		}
		if _, found := d.children[name]; found {
			return fmt.Errorf("entry already exists: %s", name)
		}
	}
	if d.children == nil {
		d.children = make(map[string]*Dirent, len(names))
	}
	for _, name = range names {
		newDirent := &Dirent{
			name:   name,
			parent: d,
		}
		d.children[name] = newDirent
	}
	return nil
}

// Add creates a new directory entry as a child of the receiver.
func (d *Dirent) Add(name string) (*Dirent, error) {
	if name == "" || name == "/" {
		return nil, errors.New("invalid name")
	}
	if _, found := d.children[name]; found {
		return nil, errors.New("entry already exists")
	}
	newDirent := &Dirent{
		name:   name,
		parent: d,
	}
	if d.children == nil {
		d.children = make(map[string]*Dirent)
	}
	d.children[name] = newDirent
	return newDirent, nil
}

// Unlink removes the dirent from its parent.
func (d *Dirent) Unlink() bool {
	if d.parent == nil {
		return false
	}
	parent := d.parent
	d.parent = nil
	ch, found := parent.children[d.name]
	if !found {
		return false
	}
	if d != ch {
		panic("parent has wrong name for child")
	}
	delete(parent.children, d.name)
	return true
}

// Move re-parents the directory entry to be a child of the given destination.
func (d *Dirent) Move(dst *Dirent) error {
	if d.name == "" {
		return errors.New("cannot move unnamed directory")
	}

	if _, found := dst.children[d.name]; found {
		return errors.New("entry already exists at destination")
	}

	d.Unlink()
	dst.children[d.name] = d
	d.parent = dst
	return nil
}

// Rename changes the name of the dirent.
// If the parent dirent already contains an entry with the new name, then an
// error is returned.
func (d *Dirent) Rename(name string) error {
	if name == "/" {
		name = ""
	}
	if d.parent != nil {
		if name == "" {
			return errors.New("non-root directory must have name")
		}
		if _, found := d.parent.children[name]; found {
			return errors.New(
				"another entry with requested name exists in parent")
		}
	}

	d.name = name
	return nil
}

// Find performs a breadth-first search for a node with the given name.
func (d *Dirent) Find(name string) *Dirent {
	nodes := queue.New()
	nodes.Add(d)
	var node, ch *Dirent
	var found bool
	for nodes.Length() > 0 {
		node = nodes.Remove().(*Dirent)
		if ch, found = node.children[name]; found {
			return ch
		}
		for _, ch = range node.children {
			nodes.Add(ch)
		}
	}
	return nil
}

// Children returns a list of names of the children of the given node.
func (d *Dirent) Children() []*Dirent {
	if len(d.children) == 0 {
		return nil
	}
	children := make([]*Dirent, 0, len(d.children))
	for _, ch := range d.children {
		children = append(children, ch)
	}
	return children
}

// List returns a sorted slice of names of the children of the receiver.
func (d *Dirent) List() []string {
	if len(d.children) == 0 {
		return nil
	}
	childNames := make([]string, 0, len(d.children))
	for name, _ := range d.children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	return childNames
}

// Path returns the slash-separated full path name of the dirent.
func (d *Dirent) Path() string {
	return d.PathDelim("/")
}

// PathDelim returns the delim-separated full path name of the dirent.
func (d *Dirent) PathDelim(delim string) string {
	parts := []string{d.name}
	p := d.parent
	for p != nil {
		parts = append(parts, p.name)
		p = p.parent
	}

	// Reverse the list so it starts at root.
	for lt, rt := 0, len(parts)-1; lt < rt; {
		parts[lt], parts[rt] = parts[rt], parts[lt]
		lt++
		rt--
	}

	// If the root dir is named same as delim, then do not show it.
	if parts[0] == delim {
		parts[0] = ""
	}

	return strings.Join(parts, delim)
}

// Tree returns a string containing the pretty-printed directory tree rooted at
// the given node.
//
// The format is similar to the UNIX/Linux "tree" utility.
// http://mama.indstate.edu/users/ice/tree/
func (d *Dirent) Tree() string {
	const (
		linkPfx = "|-- "
		contPfx = "|   "
		endlPfx = "`-- "
		blnkPfx = "    "
	)

	ss := []string{d.name}
	nodes := d.Children()

	// Reverse sort the nodes, because nodes are removed from end of list.
	sort.Sort(sort.Reverse(nodeSlice(nodes)))

	pfx := linkPfx
	var ppfx string
	var ps []string
	var newChIndex int
	var cur *Dirent

	for len(nodes) > 0 {
		cur = nodes[len(nodes)-1]
		if cur == nil {
			ps = ps[:len(ps)-1]
			ppfx = strings.Join(ps, "")
			nodes = nodes[:len(nodes)-1]
			continue
		}
		if len(nodes) == 1 || nodes[len(nodes)-2] == nil {
			pfx = endlPfx
		}
		ss = append(ss, fmt.Sprintf("%s%s%s", ppfx, pfx, cur.name))

		if len(cur.children) > 0 {
			if pfx == endlPfx {
				// Last item at level, so do not continue this level line.
				ps = append(ps, blnkPfx)
			} else {
				// More items at this level, so continue this level line.
				ps = append(ps, contPfx)
			}
			ppfx = strings.Join(ps, "")

			// Add sentinel to indicate done with depth.
			newChIndex = len(nodes)
			nodes[newChIndex-1] = nil
			// Add children for next level.
			nodes = append(nodes, cur.Children()...)
			// Reverse sort only sub-slice containing the new children.
			sort.Sort(sort.Reverse(nodeSlice(nodes[newChIndex:len(nodes)])))
			pfx = linkPfx
		} else {
			nodes = nodes[:len(nodes)-1]
		}
	}
	return strings.Join(ss, "\n")
}
