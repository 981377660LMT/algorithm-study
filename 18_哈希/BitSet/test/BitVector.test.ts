// test BitVector

import { BitVector } from '../BitVector'

describe('BitVector', () => {
  let n: number
  let nums: number[]
  let bitVector: BitVector

  // beforeEach
  beforeEach(() => {
    n = Math.floor(Math.random() * 100) + 1
    nums = Array.from({ length: n }, () => Math.floor(Math.random() * 2))
    bitVector = new BitVector(n)
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitVector.add(i)
      }
    }
    bitVector.build()
  })

  it('should support add and has', () => {
    for (let i = 0; i < n; i++) {
      if (nums[i]) {
        bitVector.add(i)
      }
    }

    for (let i = 0; i < n; i++) {
      expect(bitVector.has(i)).toBe(!!nums[i])
    }
  })

  // clear
  it('should support clear', () => {
    bitVector.clear()
    for (let i = 0; i < n; i++) {
      nums[i] = 0
    }

    for (let i = 0; i < n; i++) {
      expect(bitVector.has(i)).toBe(false)
    }
  })

  // indexOf
  it('should support kth/kthWithStart ', () => {
    const zeroIndex = nums.map((v, i) => (v === 0 ? i : -1)).filter(v => v !== -1)
    const oneIndex = nums.map((v, i) => (v === 1 ? i : -1)).filter(v => v !== -1)
    zeroIndex.forEach((v, i) => {
      expect(bitVector.kth(0, i)).toBe(v)
    })
    oneIndex.forEach((v, i) => {
      expect(bitVector.kth(1, i)).toBe(v)
    })
    const start = Math.floor(Math.random() * n)
    zeroIndex
      .filter(v => v >= start)
      .forEach((v, i) => {
        expect(bitVector.kthWithStart(0, i, start)).toBe(v)
      })
    oneIndex
      .filter(v => v >= start)
      .forEach((v, i) => {
        expect(bitVector.kthWithStart(1, i, start)).toBe(v)
      })
  })

  // count
  it('should support count', () => {
    expect(bitVector.countPrefix(0, bitVector.size)).toBe(nums.filter(v => v === 0).length)
    expect(bitVector.countPrefix(1, bitVector.size)).toBe(nums.filter(v => v === 1).length)

    // random start and end
    for (let i = 0; i < 1; i++) {
      const start = Math.floor(Math.random() * n)
      const end = Math.floor(Math.random() * (n - start)) + start
      expect(bitVector.count(1, start, end)).toBe(
        nums.slice(start, end).filter(v => v === 1).length
      )
      expect(bitVector.count(0, start, end)).toBe(
        nums.slice(start, end).filter(v => v === 0).length
      )
    }
  })
})
