package main

import (
	"bufio"
	"fmt"
	"io"
)

func SolvePersistentBookcase(reader io.Reader, writer io.Writer) {
	in := bufio.NewReader(reader)
	out := bufio.NewWriter(writer)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)

	rows := make([]*bitsetFastFlipAll, n)
	for i := range rows {
		rows[i] = newBitsetFastFlipAll(m, false)
	}

	git := NewGit(1, q)
	versions := make([]CommitID, q+1)
	versions[0] = git.Head(0)
	answers := make([]int32, q)
	oneCount := int32(0)

	add := func(row, col int32) bool {
		if rows[row].Add(col) {
			oneCount++
			return true
		}
		return false
	}

	remove := func(row, col int32) bool {
		if rows[row].Discard(col) {
			oneCount--
			return true
		}
		return false
	}

	flipRow := func(row int32) bool {
		oneCount -= rows[row].OnesCount()
		rows[row].FlipAll()
		oneCount += rows[row].OnesCount()
		return true
	}

	const master BranchID = 0
	for qi := int32(0); qi < q; qi++ {
		var op uint8
		fmt.Fscan(in, &op)

		switch op {
		case 1:
			var row, col int32
			fmt.Fscan(in, &row, &col)
			row--
			col--
			git.Commit(master, func() bool {
				return add(row, col)
			}, func() {
				remove(row, col)
			})
		case 2:
			var row, col int32
			fmt.Fscan(in, &row, &col)
			row--
			col--
			git.Commit(master, func() bool {
				return remove(row, col)
			}, func() {
				add(row, col)
			})
		case 3:
			var row int32
			fmt.Fscan(in, &row)
			row--
			git.Commit(master, func() bool {
				return flipRow(row)
			}, func() {
				flipRow(row)
			})
		case 4:
			var k int32
			fmt.Fscan(in, &k)
			git.Reset(master, versions[k])
		}

		versions[qi+1] = git.Head(master)
		queryIndex := qi
		git.Query(master, func() {
			answers[queryIndex] = oneCount
		})
	}

	git.Execute()
	for _, v := range answers {
		fmt.Fprintln(out, v)
	}
}

type bitsetFastFlipAll struct {
	flip      bool
	n         int32
	onesCount int32
	data      []uint64
}

func newBitsetFastFlipAll(n int32, filled bool) *bitsetFastFlipAll {
	data := make([]uint64, n>>6+1)
	onesCount := int32(0)
	if filled {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= int32(len(data)<<6) - n
		}
		onesCount = n
	}
	return &bitsetFastFlipAll{n: n, data: data, onesCount: onesCount}
}

func (b *bitsetFastFlipAll) FlipAll() {
	b.flip = !b.flip
	b.onesCount = b.n - b.onesCount
}

func (b *bitsetFastFlipAll) Add(i int32) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 != b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount++
	return true
}

func (b *bitsetFastFlipAll) Discard(i int32) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 == b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount--
	return true
}

func (b *bitsetFastFlipAll) OnesCount() int32 {
	return b.onesCount
}
