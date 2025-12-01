package y23d07p1

import (
	"fmt"
	"testing"
)

func TestRuns(t *testing.T) {
	cards := []Card{Ace, Ace, Queen, Queen, Queen}
	runs := getRuns(cards)

	fmt.Printf("%+v\n", runs)
}
