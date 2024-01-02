package sets

import (
	"fmt"
	"math/bits"
	"strings"
)

type BitSet16 uint16

func (b *BitSet16) Set(e uint8) BitSet16 {
	if e > 15 {
		panic(fmt.Errorf("too big: %d", e))
	}
	return (*b) | (1 << e)
}

func (b *BitSet16) Has(e uint8) bool {
	return ((*b) & (1 << e)) != 0
}

func (b *BitSet16) Key() uint16 {
	return uint16(*b)
}

func (b *BitSet16) Len() int {
	return bits.OnesCount16(uint16(*b))
}

func (b *BitSet16) String() string {
	var buf strings.Builder
	buf.WriteString("{")
	for i := uint8(0); i < 32; i++ {
		if b.Has(i) {
			fmt.Fprintf(&buf, " %d", i)
		}
	}
	buf.WriteString(" }")
	return buf.String()
}
