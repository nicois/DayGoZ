package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
)

type Item interface {
	Use(p *Person) bool // returns if it is consumed
	Combine(other Item) Item
	GetSize() (int, int) // width, height
	Icon() rune
}

type Holder interface {
	Insert(Item) bool
	Remove(Item) error
	ListContents() []Item
	String() string
}

type backpack struct {
	mutex      *sync.Mutex
	width      int
	height     int
	contents   []Item
	contentMap map[int]Item
}

func (self *backpack) String() string {
	parts := make([]string, 0, 20)
	tl := make(map[Item]bool)
	parts = append(parts, "Pack contents:")
	for index := 0; index < self.width*self.height; index++ {
		if index%(self.width) == 0 {
			parts = append(parts, "\n")
		}
		item, found := self.contentMap[index]
		if !found {
			parts = append(parts, "..")
		} else {
			_, found := tl[item]
			if found {
				parts = append(parts, "XX")
			} else {
				log.Printf("len of %v is %v", string(item.Icon()), len(string(item.Icon())))
				parts = append(parts, string(item.Icon())+" ")
				tl[item] = true
			}
		}
	}
	return strings.Join(append(parts, "\n"), "")
}

func CreateFieldBackpack() *backpack {
	width := 8
	height := 10
	return &backpack{width: width, height: height, mutex: new(sync.Mutex), contentMap: make(map[int]Item, width*height)}
}

func CreateSchoolBackpack() *backpack {
	width := 4
	height := 6
	return &backpack{width: width, height: height, mutex: new(sync.Mutex), contentMap: make(map[int]Item, width*height)}
}

/*
    0  1  2  3
    4  5  6  7
    8  9 10 11
   12 13 14 15
   16 17 18 19
   20 21 22 23
   index(1,3) == 7
   index(1,2) == 6
*/
func (self *backpack) index(x, y int) int {
	return self.width*y + x
}

func (self *backpack) Insert(i Item) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	width, height := i.GetSize()
nextrow:
	for y := 0; y < self.height; y++ {
		log.Printf("checking y=%v\n", y)
	nextpos:
		for x := 0; x < self.width; x++ {
			index := self.index(x, y)
			log.Printf("Checking at x:%v, y:%v, index:%v.\n", x, y, index)
			for a := 0; a < width; a++ {
				for b := 0; b < height; b++ {
					if x+a > self.width || y+b > self.height {
						log.Printf("%v/%v is no good: too far\n", x+a, y+b)
						continue nextrow
					}
					if other, used := self.contentMap[index+b*self.width+a]; used {
						log.Printf("%v/%v is no good: used by %v\n", x+a, y+b, other)
						continue nextpos
					}
				}
			}
			log.Printf("Inserting %v at %v/%v\n", i, x, y)
			self.contents = append(self.contents, i)
			for a := 0; a < width; a++ {
				for b := 0; b < height; b++ {
					self.contentMap[index+b*self.width+a] = i
					log.Printf("Inserted at index %v\n", index+b*self.width+a)
				}
			}
			return true
		}
	}
	return false
}

func (self *backpack) Remove(i Item) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return nil
}

func (self *backpack) ListContents() []Item {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.contents
}

func Demonstration(p *Person) {
	var item Item
	if rand.Intn(10) < 5 {
		item = Peaches()
	} else {
		item = Chocolate()
	}

	p.Log(fmt.Sprintf("Added %v to pack: %v", item, p.pack.Insert(item)))
}
