import { SegmentTree01 } from '../SegmentTree01'

const INF = 2e15

describe('SegmentTree01', () => {
  const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
  const n = randint(100, 1000)
  const nums = Array.from({ length: n }, () => randint(0, 1))
  const seg = new SegmentTree01(nums)

  it('should support onesCount', () => {
    for (let i = 0; i < 10; i++) {
      const l = randint(1, n)
      const r = randint(l, n)
      const ones = nums.slice(l - 1, r).reduce((a, b) => a + b, 0)
      expect(seg.onesCount(l, r)).toBe(ones)
    }
  })

  it('should support flip', () => {
    for (let i = 0; i < 10; i++) {
      const l = randint(1, n)
      const r = randint(l, n)
      seg.flip(l, r)
      for (let j = l - 1; j < r; j++) {
        nums[j] ^= 1
      }
      for (let j = 1; j <= n; j++) {
        expect(seg.onesCount(j, j)).toBe(nums[j - 1])
      }
    }
  })

  it('should support lastIndexOf', () => {
    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const searchPos = randint(1, n)
      const res = seg.lastIndexOf(digit, searchPos)
      if (res === -1) {
        expect(nums.lastIndexOf(digit, searchPos - 1)).toBe(-1)
      } else {
        expect(res).toBe(nums.lastIndexOf(digit, searchPos - 1) + 1)
      }
    }
  })

  it('should support indexOf', () => {
    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const searchPos = randint(1, n)
      const res = seg.indexOf(digit, searchPos)
      if (res === -1) {
        expect(nums.indexOf(digit, searchPos - 1)).toBe(-1)
      } else {
        expect(res).toBe(nums.indexOf(digit, searchPos - 1) + 1)
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
            return i + 1
          }
        }
      }
      return -1
    }

    for (let i = 0; i < 10; i++) {
      const digit = randint(0, 1) as 0 | 1
      const k = randint(1, n)
      const res1 = seg.kth(digit, k)
      const res2 = countKth(nums, digit, k)
      expect(res1).toBe(res2)
    }
  })
})
