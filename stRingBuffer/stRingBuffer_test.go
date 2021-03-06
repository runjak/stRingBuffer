package stRingBuffer

import (
	"fmt"
	"strconv"
	"testing"
)

func TestInfo(t *testing.T) {
	//Starting with an empty buffer:
	srb := NewStRingBuffer(5)
	if !srb.Empty() {
		t.Errorf("StRingBuffer was not empty after NewStRingBuffer(x).")
	}
	//Pushing some elements:
	srb.Push("1", "2", "3", "4", "5")
	l := srb.Length()
	if l != 5 {
		t.Errorf("StRingBuffer didn't grow to expected Length(), but to %i.", l)
	}
	//Buffer should be full:
	if !srb.Full() {
		t.Errorf("StRingBuffer isn't Full(), but we expected it to be.")
	}
	//String output:
	want := "StRingBuffer{start: 0, end: 0, lines: [1 2 3 4 5], length: 5}"
	get := fmt.Sprintf("%s", srb) // Mainly so that 'fmt' can be imported.
	if get != want {
		t.Errorf("StRingBuffer String() test failed with get='%s'.\n", get)
	}
}

func TestAddTake(t *testing.T) {
	srb := NewStRingBuffer(3)
	//Testing Push/Pop:
	srb.Push("1", "2", "3", "4")
	s := ""
	for !srb.Empty() {
		s += srb.Pop()
	}
	if s != "432" {
		t.Errorf("StRingBuffer Pop() didn't reverse Push(), and gave '%s'.", s)
	}
	//Testing Un-/Shift:
	srb.Unshift("4", "3", "2", "1")
	s = ""
	for !srb.Empty() {
		s += srb.Shift()
	}
	if s != "123" {
		t.Errorf("StRingBuffer Shift() didn't reverse Unshift(), and gave '%s'.", s)
	}
	//Testing Push/Shift:
	srb.Push("1", "2", "3", "4")
	s = ""
	for !srb.Empty() {
		s += srb.Shift()
	}
	if s != "234" {
		t.Errorf("StRingBuffer Shift() didn't hold the last Push() in order, but '%s'.", s)
	}
	//Testing Unshift/Pop:
	srb.Unshift("4", "3", "2", "1")
	s = ""
	for !srb.Empty() {
		s += srb.Pop()
	}
	if s != "321" {
		t.Errorf("StRingBuffer Pop() didn't hold the last Unshift() in order, but '%s'.", s)
	}
}

func TestMapEach(t *testing.T) {
	srb := NewStRingBuffer(3)
	srb.Push(".", ".", ".")
	c := 0
	f := func(s string) string {
		c++
		return s + strconv.Itoa(c)
	}
	srb.Map(f).MapR(f)
	sum := ""
	g := func(s string) {
		sum += s
	}
	srb.Each(g).EachR(g)
	if sum != ".16.25.34.34.25.16" {
		t.Errorf("StRingBuffer MapEach test didn't hold the expected result but '%s'.\n", sum)
	}
}

func TestUnsliceSlice(t *testing.T) {
	//A helper function to compare two slices of strings:
	test := func(xs []string, ys []string) {
		if len(xs) != len(ys) {
			t.Errorf("UnsliceSlice didn't preserve length of the slice.")
			return
		}
		for i, x := range xs {
			if ys[i] != x {
				t.Errorf("UnsliceSlice changed at least one element: %v != %v", xs, ys)
				return
			}
		}
	}
	//Our test data:
	xs := []string{"a", "b", "c"}
	//Does Unslice->Slice hold:
	test(xs, Unslice(xs).Slice())
	//Pushing into a new buffer, does Slice->Unslice->Slice hold:
	srb := NewStRingBuffer(3)
	srb = Unslice(srb.Push(xs...).Slice())
	test(xs, srb.Slice())
	//Is Length correct?
	if srb.Length() != len(xs) {
		t.Errorf("Slice->Unslice didn't preserver length.")
	}
	//Is Shift as expected?
	if srb.Shift() != "a" {
		t.Errorf("First element wasn't as expected.")
	}
	//Is Pop as expected?
	if l := srb.Pop(); l != "c" {
		t.Errorf("Last element wasn't as expected, \"%s\" != \"c\"", l)
	}
}
