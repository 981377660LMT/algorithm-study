// test BitSet

import { BitSet } from '../BitSet'

describe('BitSet', () => {
  let n: number
  let nums: number[]
  let bitSet: BitSet

  // beforeEach
  beforeEach(() => {
    n = Math.floor(Math.random() * 1000) + 1
    nums = Array.from({ length: n }, () => Math.floor(Math.random() * 2))
    bitSet = new BitSet(n)
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitSet.add(i)
      }
    }
  })

  it('should support add and has', () => {
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitSet.add(i)
      }
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // addRange
  it('should support addRange', () => {
    const start = Math.floor(Math.random() * n)
    const end = Math.floor(Math.random() * (n - start)) + start
    bitSet.addRange(start, end)
    for (let i = start; i < end; i++) {
      nums[i] = 1
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // discard
  it('should support discard', () => {
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitSet.discard(i)
        nums[i] = 0
      }
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(false)
    }
  })

  // discardRange
  it('should support discardRange', () => {
    const start = Math.floor(Math.random() * n)
    const end = Math.floor(Math.random() * (n - start)) + start
    bitSet.discardRange(start, end)
    for (let i = start; i < end; i++) {
      nums[i] = 0
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // flip
  it('should support flip', () => {
    for (let i = 0; i < n; i++) {
      bitSet.flip(i)
      nums[i] ^= 1
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // flipRange
  it('should support flipRange', () => {
    const start = Math.floor(Math.random() * n)
    const end = Math.floor(Math.random() * (n - start)) + start
    bitSet.flipRange(start, end)
    for (let i = start; i < end; i++) {
      nums[i] ^= 1
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // clear
  it('should support clear', () => {
    bitSet.clear()
    for (let i = 0; i < n; i++) {
      nums[i] = 0
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(false)
    }
  })

  it('should support allOne and allZero', () => {
    // random add
    for (let i = 0; i < n; i++) {
      if (Math.random() > 0.5) {
        bitSet.add(i)
        nums[i] = 1
      }
    }

    // random start and end
    const start = Math.floor(Math.random() * n)
    const end = Math.floor(Math.random() * (n - start)) + start

    // check
    for (let i = start; i < end; i++) {
      expect(bitSet.allOne(start, end)).toBe(nums.slice(start, end).every(v => v === 1))
    }

    // random start and end
    const start2 = Math.floor(Math.random() * n)
    const end2 = Math.floor(Math.random() * (n - start2)) + start2

    // check
    for (let i = start2; i < end2; i++) {
      expect(bitSet.allZero(start2, end2)).toBe(nums.slice(start2, end2).every(v => v === 0))
    }
  })

  // indexOfOne
  it('should support indexOfOne ', () => {
    for (let i = 0; i < n; i++) {
      expect(bitSet.indexOfOne(i)).toBe(nums.indexOf(1, i))
    }
  })

  // lastIndexOfOne
  it('should support lastIndexOfOne ', () => {
    expect(bitSet.bitLength() - 1).toBe(nums.lastIndexOf(1))
  })

  // onesCount
  it('should support onesCount', () => {
    expect(bitSet.onesCount()).toBe(nums.reduce((a, b) => a + b, 0))

    // random start and end
    for (let i = 0; i < 1; i++) {
      const start = Math.floor(Math.random() * n)
      const end = Math.floor(Math.random() * (n - start)) + start
      expect(bitSet.onesCount(start, end)).toBe(nums.slice(start, end).reduce((a, b) => a + b, 0))
    }
  })

  // equals and copy
  it('should support equals and copy', () => {
    const bitSet2 = bitSet.copy()
    expect(bitSet.equals(bitSet2)).toBe(true)
  })

  // copy with newSize (copy and resize)
  it('should support copy with newSize', () => {
    const newSize = Math.floor(Math.random() * n * 2)
    const bitSet2 = bitSet.copy(newSize)
    for (let i = 0; i < Math.min(n, newSize); i++) {
      expect(bitSet2.has(i)).toBe(bitSet.has(i))
    }
    for (let i = n; i < newSize; i++) {
      expect(bitSet2.has(i)).toBe(false)
    }
  })

  // iOr
  it('should support iOr', () => {
    const bitSet2 = new BitSet(n)
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitSet2.add(i)
      }
    }

    bitSet.ior(bitSet2)
    for (let i = 0; i < n; i++) {
      nums[i] = nums[i] || nums[i]
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // iAnd
  it('should support iAnd', () => {
    const bitSet2 = new BitSet(n)
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitSet2.add(i)
      }
    }

    bitSet.iand(bitSet2)
    for (let i = 0; i < n; i++) {
      nums[i] = nums[i] && nums[i]
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // isSubset/isSuperset
  it('should support isSubset/isSuperset', () => {
    const bitSet2 = bitSet.copy()
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        // random fail
        if (Math.random() > 0.5) {
          bitSet2.discard(i)
        }
      }
    }

    expect(bitSet2.isSubset(bitSet)).toBe(true)
    expect(bitSet.isSuperset(bitSet2)).toBe(true)
  })

  // fill
  it('should support fill', () => {
    const bitSet2 = new BitSet(n)
    const zeroOrOne = Math.random() > 0.5 ? 1 : 0
    bitSet2.fill(zeroOrOne)
    for (let i = 0; i < n; i++) {
      nums[i] = zeroOrOne
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet2.has(i)).toBe(!!nums[i])
    }
  })

  // slice
  it('should support slice', () => {
    for (let _ = 0; _ < 1000; _++) {
      const start = Math.floor(Math.random() * n)
      const end = Math.floor(Math.random() * (n - start)) + start
      const bitSet2 = bitSet.slice(start, end)
      for (let i = start; i < end; i++) {
        expect(bitSet2.has(i - start)).toBe(!!nums[i])
      }
    }
  })

  // set
  it('should support set', () => {
    for (let i = 0; i < 10; i++) {
      const m = (n * Math.random()) | 0
      const newBitSet = new BitSet(m)
      const newNums: number[] = Array(m).fill(0)
      for (let i = 0; i < m; i++) {
        const cur = Math.random() > 0.5 ? 1 : 0
        if (cur) {
          newBitSet.add(i)
          newNums[i] = 1
        }
      }

      const offset = ((n - m) * Math.random()) | 0
      bitSet.set(newBitSet, offset)
      for (let i = 0; i < m; i++) {
        nums[i + offset] = newNums[i]
      }

      for (let i = 0; i < n; i++) {
        expect(bitSet.has(i)).toBe(!!nums[i])
      }
    }
  })

  // iorRange(start, end, other)
  it('should support iorRange', () => {
    const m = (n * Math.random()) | 0
    const newBitSet = new BitSet(m)
    const newNums: number[] = Array(m).fill(0)
    for (let i = 0; i < m; i++) {
      const cur = Math.random() > 0.5 ? 1 : 0
      if (cur) {
        newBitSet.add(i)
        newNums[i] = 1
      }
    }

    const offset = ((n - m) * Math.random()) | 0
    bitSet.iorRange(offset, offset + m, newBitSet)
    for (let i = 0; i < m; i++) {
      nums[i + offset] = nums[i + offset] || newNums[i]
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // iandRange(start, end, other)
  it('should support iandRange', () => {
    const m = (n * Math.random()) | 0
    const newBitSet = new BitSet(m)
    const newNums: number[] = Array(m).fill(0)
    for (let i = 0; i < m; i++) {
      const cur = Math.random() > 0.5 ? 1 : 0
      if (cur) {
        newBitSet.add(i)
        newNums[i] = 1
      }
    }

    const offset = ((n - m) * Math.random()) | 0
    bitSet.iandRange(offset, offset + m, newBitSet)
    for (let i = 0; i < m; i++) {
      nums[i + offset] = nums[i + offset] && newNums[i]
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // ixorRange(start, end, other)
  it('should support ixorRange', () => {
    const m = (n * Math.random()) | 0
    const newBitSet = new BitSet(m)
    const newNums: number[] = Array(m).fill(0)
    for (let i = 0; i < m; i++) {
      const cur = Math.random() > 0.5 ? 1 : 0
      if (cur) {
        newBitSet.add(i)
        newNums[i] = 1
      }
    }

    const offset = ((n - m) * Math.random()) | 0
    bitSet.ixorRange(offset, offset + m, newBitSet)
    for (let i = 0; i < m; i++) {
      nums[i + offset] = +(nums[i + offset] !== newNums[i])
    }

    for (let i = 0; i < n; i++) {
      expect(bitSet.has(i)).toBe(!!nums[i])
    }
  })

  // next
  it('should support next', () => {
    for (let i = 0; i < n; i++) {
      const index1 = bitSet.next(i)
      const index2 = nums.indexOf(1, i)
      if (index2 === -1) {
        expect(index1).toBe(bitSet.size)
      } else {
        expect(index1).toBe(index2)
      }
    }
  })

  // prev
  it('should support prev', () => {
    for (let i = n - 1; i > 0; i--) {
      const index1 = bitSet.prev(i)
      const index2 = nums.lastIndexOf(1, i)
      if (index2 === -1) {
        expect(index1).toBe(-1)
      } else {
        expect(index1).toBe(index2)
      }
    }
  })

  // lsh
  it('should support lsh', () => {
    for (let i = 0; i < n; i++) {
      const shift = Math.floor(Math.random() * n)
      const set1 = bitSet.toSet()
      bitSet.lsh(shift)
      const set2 = bitSet.toSet()
      for (const v of set1) {
        if (v + shift <= n) {
          expect(set2.has(v + shift)).toBe(true)
        }
      }
    }
  })

  // rsh
  it('should support rsh', () => {
    for (let i = 0; i < n; i++) {
      const shift = Math.floor(Math.random() * n)
      const set1 = bitSet.toSet()
      bitSet.rsh(shift)
      const set2 = bitSet.toSet()
      for (const v of set1) {
        if (v - shift >= 0) {
          expect(set2.has(v - shift)).toBe(true)
        }
      }
    }
  })
})
