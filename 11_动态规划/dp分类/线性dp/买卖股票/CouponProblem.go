package main

func main() {
	price := []int{3, 5, 8, 2, 6}
	discount := []int{2, 2, 4, 1, 3}
	numCoupon := int32(2)
	money := 20
	res := CouponProblem(price, discount, numCoupon, money)
	for i, choice := range res {
		if choice == NotBuy {
			continue
		}
		if choice == BuyWithoutCoupon {
			println("BuyWithoutCoupon:", price[i])
		} else {
			println("BuyWithCoupon:", discount[i])
		}
	}
}

type Choice uint8

const (
	NotBuy Choice = iota
	BuyWithoutCoupon
	BuyWithCoupon
)

// 商品问题(股票问题).
// price 表示每份商品的价格, discount 表示每份商品使用优惠券后的价格.
// numCoupon 表示可以使用的优惠券数量, money 表示总金额.
func CouponProblem(price []int, discount []int, numCoupon int32, money int) []Choice {
	n := int32(len(price))
	items := make([]*Item, n)
	for i := int32(0); i < n; i++ {
		items[i] = &Item{
			price:    price[i],
			discount: discount[i],
			profit:   price[i] - discount[i],
			index:    i,
		}
	}
	res := make([]Choice, n)
	discountGroup := NewHeap(func(a, b *Item) bool { return a.profit < b.profit }, nil)
	withoutDiscount := NewHeap(func(a, b *Item) bool { return a.price < b.price }, items)
	withDiscount := NewHeap(func(a, b *Item) bool { return a.discount < b.discount }, items)
	for withDiscount.Len() > 0 && withoutDiscount.Len() > 0 {
		if res[withDiscount.Top().index] != NotBuy {
			withDiscount.Pop()
			continue
		}
		if res[withoutDiscount.Top().index] != NotBuy {
			withoutDiscount.Pop()
			continue
		}
		discountTop := 0
		if numCoupon == 0 {
			discountTop = money + 1
		} else if discountGroup.Len() == int(numCoupon) {
			discountTop = discountGroup.Top().profit
		}
		if discountTop+withDiscount.Top().discount <= withoutDiscount.Top().price {
			top := withDiscount.Pop()
			cost := discountTop + top.discount
			if cost > money {
				break
			}
			money -= cost
			res[top.index] = BuyWithCoupon
			if discountGroup.Len() == int(numCoupon) {
				res[discountGroup.Pop().index] = BuyWithoutCoupon
			}
			discountGroup.Push(top)
		} else {
			top := withoutDiscount.Pop()
			cost := top.price
			if cost > money {
				break
			}
			money -= cost
			res[top.index] = BuyWithoutCoupon
		}
	}
	return res
}

type Item struct {
	price, discount, profit int
	index                   int32
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
