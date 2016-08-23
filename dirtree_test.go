package dirtree

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	dt := New("")
	fmt.Println(dt.Tree())
	s1, _ := dt.Add("A")
	fmt.Println(dt.Tree())
	s2, _ := dt.Add("B")
	fmt.Println(dt.Tree())
	s1.Add("Ax")
	fmt.Println(dt.Tree())
	s3, _ := s1.Add("Ay")
	fmt.Println(dt.Tree())
	s2.Add("Bx")
	fmt.Println(dt.Tree())
	s2.Add("By")
	fmt.Println(dt.Tree())
	s3.Add("Ay-1")
	fmt.Println(dt.Tree())
	s4, _ := s3.Add("Ay-2")
	fmt.Println(dt.Tree())
	fmt.Println()
	fmt.Println(s3.Tree())
	fmt.Println()
	fmt.Println(s4.Path())
}

func TestTree(t *testing.T) {
	dt := New(".")
	if dt.Tree() != "." {
		t.Errorf("Expected '.' with empty tree")
	}

	s1, _ := dt.Add("A")
	if dt.Tree() != ".\n`-- A" {
		t.Errorf("Incorrect output with single entry")
	}

	s2, _ := dt.Add("B")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"`-- B" {
		t.Errorf("Incorrect output for two entries at same level")
	}

	s1.Add("Ax")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   `-- Ax\n"+
			"`-- B" {
		t.Errorf("Incorrect output for one subentry")
	}

	s3, _ := s1.Add("Ay")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   |-- Ax\n"+
			"|   `-- Ay\n"+
			"`-- B" {
		t.Errorf("Incorrect output for two sub entries")
	}

	s2.Add("Bx")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   |-- Ax\n"+
			"|   `-- Ay\n"+
			"`-- B\n"+
			"    `-- Bx" {
		t.Errorf("Incorrect output for multiple subs under different entries")
	}

	s2.Add("By")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   |-- Ax\n"+
			"|   `-- Ay\n"+
			"`-- B\n"+
			"    |-- Bx\n"+
			"    `-- By" {
		t.Errorf("Incorrect output for multiple two subs under all entries")
	}

	s3.Add("Ay-1")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   |-- Ax\n"+
			"|   `-- Ay\n"+
			"|       `-- Ay-1\n"+
			"`-- B\n"+
			"    |-- Bx\n"+
			"    `-- By" {
		t.Errorf("Incorrect output for sub at depth three")
	}

	s3.Add("Ay-2")
	if dt.Tree() !=
		".\n"+
			"|-- A\n"+
			"|   |-- Ax\n"+
			"|   `-- Ay\n"+
			"|       |-- Ay-1\n"+
			"|       `-- Ay-2\n"+
			"`-- B\n"+
			"    |-- Bx\n"+
			"    `-- By" {
		t.Errorf("Incorrect output for multiple subs at depth three")
	}

	d1 := dt.Tree()
	d2 := dt.Tree()
	if d1 != d2 {
		t.Errorf("Output not stable across multiple calls to String()")
	}
}

func TestChild(t *testing.T) {
	dt := New("")
	s1, _ := dt.Add("A")
	s2, _ := dt.Add("B")
	s11, _ := s1.Add("Ax")
	s1.Make("Ay", "Az")
	s2.Add("Bx")
	s22, _ := s2.Add("By")
	if dt.Child("A") != s1 {
		t.Fatalf("returned wrong child")
	}
	if dt.Child("B") != s2 {
		t.Fatalf("returned wrong child")
	}
	if s1.Child("Ax") != s11 {
		t.Fatalf("returned wrong child")
	}
	if s2.Child("By") != s22 {
		t.Fatalf("returned wrong child")
	}
	if s2.Child("NotHere") != nil {
		t.Fatalf("expected nil")
	}
}

func TestList(t *testing.T) {
	dt := New("")
	dt.Make("Z", "X", "Y")
	dt.Add("A")
	dt.Add("C")
	dt.Add("B")
	names := dt.List()
	if len(names) != 6 {
		t.Fatal("list has incorrect length")
	}
	expected := []string{"A", "B", "C", "X", "Y", "Z"}
	for i, name := range names {
		if name != expected[i] {
			t.Fatal("incorrect children names list")
		}
	}
}

func TestChildren(t *testing.T) {
	dt := New("")
	z, _ := dt.Add("Z")
	x, _ := dt.Add("X")
	y, _ := dt.Add("Y")
	a, _ := dt.Add("A")
	c, _ := dt.Add("C")
	b, _ := dt.Add("B")
	children := dt.Children()
	if len(children) != 6 {
		t.Fatal("children has incorrect length")
	}
	Sort(children)
	expected := []*Dirent{a, b, c, x, y, z}
	for i, ch := range children {
		if ch != expected[i] {
			t.Fatalf("expected child %s got %s", expected[i], ch)
		}
	}

	SortReverse(children)
	expected = []*Dirent{z, y, x, c, b, a}
	for i, ch := range children {
		if ch != expected[i] {
			t.Fatalf("expected child %s got %s", expected[i], ch)
		}
	}
}

func TestFind(t *testing.T) {
	dt := New("")
	s1, _ := dt.Add("A")
	s2, _ := dt.Add("B")
	s11, _ := s1.Add("Ax")
	s1.Make("Ay", "Az")
	s2.Add("Bx")
	s2.Add("By")
	s111, _ := s11.Add("Ax-1")
	found := dt.Find("Ax-1")
	if found == nil {
		t.Fatal("node not found")
	}
	if found != s111 {
		t.Fatal("found wrong node")
	}
	fmt.Println("Found:", found.Path())
}
