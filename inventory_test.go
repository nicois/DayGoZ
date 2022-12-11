package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestRenderPack(t *testing.T) {
	width := 4
	height := 6
	pack := &backpack{width: width, height: height, mutex: new(sync.Mutex), contentMap: make(map[int]Item, width*height)}
	pack.Insert(Peaches())
	pack.Insert(Peaches())
	pack.Insert(Chocolate())
	pack.Insert(Chocolate())
	rendered := fmt.Sprint(pack)
	t.Fatal(rendered)
}

func TestInsert(t *testing.T) {
	width := 4
	height := 6
	pack := &backpack{width: width, height: height, mutex: new(sync.Mutex), contentMap: make(map[int]Item, width*height)}
	peach1 := Peaches()
	pack.Insert(peach1)
	t.Logf("----------------\n")
	peach2 := Peaches()
	pack.Insert(peach2)
	t.Logf("----------------\n")

	nSlots := make(map[Item]int)
	for x := 0; x < 4; x++ {
		for y := 0; y < 6; y++ {
			index := y*4 + x
			if item, found := pack.contentMap[index]; found {
				t.Logf("Found %v at %v", item, index)
				nSlots[item]++
			}
		}
	}
	if len(nSlots) != 2 {
		t.Fatal(len(nSlots))
	}
	if nSlots[peach1] != 4 {
		t.Fatal(peach1)
	}
	if nSlots[peach2] != 4 {
		t.Fatal(peach2)
	}

	peach3 := Peaches()
	pack.Insert(peach3)
	t.Logf("----------------\n")

	nSlots = make(map[Item]int)
	for x := 0; x < 4; x++ {
		for y := 0; y < 6; y++ {
			index := x*6 + y
			if item, found := pack.contentMap[index]; found {
				t.Logf("Found %v at %v", item, index)
				nSlots[item]++
			}
		}
	}
	if len(nSlots) != 3 {
		t.Fatal(len(nSlots))
	}
	if nSlots[peach1] != 4 {
		t.Fatal(peach1)
	}
	if nSlots[peach2] != 4 {
		t.Fatal(peach2)
	}
	if nSlots[peach3] != 4 {
		t.Fatal(peach3)
	}

	// A total of 6 peaches should fix in a 4/6 pack
	for i := 0; i < 3; i++ {
		if !pack.Insert(Peaches()) {
			t.Fatal("Could not insert peaches")
		}
	}
	// Next attempted insertion should fail
	if pack.Insert(Peaches()) {
		t.Fatal("Could insert 7th peaches")
	}
}
