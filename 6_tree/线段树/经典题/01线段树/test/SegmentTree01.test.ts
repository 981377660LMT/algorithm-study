import { SegmentTree01 } from '../SegmentTree01'

const INF = 2e15

describe('SegmentTree01', () => {
  const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
  const n = randint(100, 500)
  const nums = Array.from({ length: n }, () => randint(0, 1))
  const seg = new SegmentTree01(nums)

  // constructor bits/length
  it('should support constructor bits/length', () => {
    const n = randint(100, 500)
    const seg1 = new SegmentTree01(Array(n).fill(0))
    const seg2 = new SegmentTree01(n)
    expect(seg1.toString()).toBe(seg2.toString())
  })

  it('should support onesCount', () => {
    for (let i = 0; i < 10; i++) {
      const l = randint(1, n + 10)
      const r = randint(l, n + 10)
      const ones = nums.slice(l, r).reduce((a, b) => a + b, 0)
      expect(seg.onesCount(l, r)).toBe(ones)
    }
  })

  it('should support flip', () => {
    for (let i = 0; i < 10; i++) {
      const l = randint(1, n + 10)
      const r = randint(l, n + 10)
      seg.flip(l, r)
      for (let j = l; j < r; j++) {
        nums[j] ^= 1
      }
      for (let j = 1; j <= n; j++) {
        expect(seg.onesCount(j - 1, j)).toBe(nums[j - 1])
      }
    }
  })

  it('should support lastIndexOf', () => {
    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const searchPos = randint(0, n)
      const res = seg.lastIndexOf(digit, searchPos)
      if (res === -1) {
        expect(nums.lastIndexOf(digit, searchPos)).toBe(-1)
      } else {
        expect(res).toBe(nums.lastIndexOf(digit, searchPos))
      }
    }
  })

  it('should support indexOf', () => {
    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const searchPos = randint(0, n)
      const res = seg.indexOf(digit, searchPos)
      if (res === -1) {
        expect(nums.indexOf(digit, searchPos)).toBe(-1)
      } else {
        expect(res).toBe(nums.indexOf(digit, searchPos))
      }
    }
  })

  it('should support kth', () => {
    const countKth = (nums: number[], digit: 0 | 1, k: number) => {
      let count = 0
      for (let i = 0; i < nums.length; i++) {
        if (nums[i] === digit) {
          count++
          if (count === k) {
            return i
          }
        }
      }
      return -1
    }

    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const k = randint(1, n)
      const res1 = seg.kth(digit, k)
      let res2 = countKth(nums, digit, k)
      expect(res1).toBe(res2)
    }
  })
})
