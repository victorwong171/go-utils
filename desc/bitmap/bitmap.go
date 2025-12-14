package bitmap

type bitMap struct {
	vmax  uint
	graph []byte
}

func NewBitMap(max uint) *bitMap {
	var defMax uint = 8192
	if max > 0 {
		defMax = max
	}
	size := (max + 7) >> 3
	return &bitMap{
		vmax:  defMax,
		graph: make([]byte, size, size),
	}
}

func (b *bitMap) Set(num uint) {
	if num > b.vmax {
		b.vmax += 1024
		if num > b.vmax {
			b.vmax = num
		}
		inc := int(b.vmax+7)>>3 - len(b.graph)
		if inc > 0 {
			b.graph = append(b.graph, make([]byte, inc, inc)...)
		}
	}
	b.graph[num>>3] |= 1 << (num % 8)
}

func (b *bitMap) Del(num uint) {
	if num <= b.vmax {
		b.graph[num>>3] &^= 1 << (num % 8)
	}
}

func (b *bitMap) Check(num uint) bool {
	if num > b.vmax {
		return false
	}
	return b.graph[num>>3]&(1<<(num%8)) != 0
}
