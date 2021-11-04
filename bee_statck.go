package beehive

import (
	"time"
)

type beeStack struct {
	items  []*bee
	expiry []*bee
	size   int
}

func newBeeStack(size int) *beeStack {
	return &beeStack{
		items: make([]*bee, 0, size),
		size:  size,
	}
}

func (bs *beeStack) len() int {
	return len(bs.items)
}

func (bs *beeStack) isEmpty() bool {
	return bs.len() == 0
}

func (bs *beeStack) insert(w *bee) error {
	bs.items = append(bs.items, w)
	return nil
}

func (bs *beeStack) detach() *bee {
	l := bs.len()
	if l == 0 {
		return nil
	}

	w := bs.items[l-1]
	bs.items[l-1] = nil
	bs.items = bs.items[:l-1]

	return w
}

func (bs *beeStack) retrieve(duration time.Duration) []*bee {
	n := bs.len()
	if n == 0 {
		return nil
	}

	expiryTime := time.Now().Add(-duration)
	index := bs.search(0, n-1, expiryTime)

	bs.expiry = bs.expiry[:0]
	if index != -1 {
		bs.expiry = append(bs.expiry, bs.items[:index+1]...)
		m := copy(bs.items, bs.items[index+1:])
		for i := m; i < n; i++ {
			bs.items[i] = nil
		}
		bs.items = bs.items[:m]
	}
	return bs.expiry
}

func (bs *beeStack) search(l, r int, expiryTime time.Time) int {
	var mid int
	for l < r {
		mid = (l + r) / 2
		if expiryTime.Before(bs.items[mid].recyleTime) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	return r
}

func (bs *beeStack) reset() {
	for i := 0; i < bs.len(); i++ {
		bs.items[i].task <- nil
		bs.items[i] = nil
	}
	bs.items = bs.items[:0]
}
