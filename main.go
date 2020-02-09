package main

import (
	crand "crypto/rand"
	"fmt"
	rand "math/rand"

	"encoding/binary"
	"log"
)

func main() {

	winOnePick := 0
	winSecondPick := 0

	count := 10000000
	for index := 0; index < count; index++ {
		doors := initDoors()
		firstPick := pickFirst()

		if doors[pickFirst()] == true {
			winOnePick++
		}

		var pickedDoor int

		for index, door := range doors {
			if index == firstPick {
				// Ignore picked door, because we already picked it and we are going to switch
				continue
			}
			if doors[firstPick] == true && door != true {
				// If our pick was a win, it does not matter which pick we make, because it's a lose.
				// The host will open a "bad" door, and the only choice left is another bad door.
				// Hence the 2/3 chance to win, we only lose if we actually picked the right door in the first place
				pickedDoor = index
			} else {
				if door == true {
					// If our picked door was a "bad" door, the host will now open another "bad" door.
					// Obviously we don't pick a "bad" door, which causes a default win.
					pickedDoor = index
				}
			}
		}

		if doors[pickedDoor] == true {
			winSecondPick++
		}

	}

	fmt.Println("Runs:", count)
	fmt.Printf("Wins with the first pick only: %d ( %0.2f %% ) \n", winOnePick, percentage(winOnePick, count))
	fmt.Printf("Wins with the second pick: %d ( %0.2f %% ) \n", winSecondPick, percentage(winSecondPick, count))

}

func initDoors() [3]bool {

	doors := [3]bool{false, false, false}

	var src cryptoSource
	rnd := rand.New(src)
	doors[rnd.Intn(3)] = true

	return doors
}

func pickFirst() int {
	var src cryptoSource
	rnd := rand.New(src)
	return rnd.Intn(3)
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func percentage(first, second int) (delta float64) {
	return float64(first) / float64(second) * 100

}
