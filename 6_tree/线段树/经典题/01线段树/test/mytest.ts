import { SegmentTree01 } from '../SegmentTree01'

const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
const n = randint(100, 500)
const nums = Array.from({ length: n }, () => randint(0, 1))
const seg = new SegmentTree01(nums)
for (let i = 0; i < 1000; i++) {
  const l = randint(1, n + 10)
  const r = randint(l, n + 10)
  seg.flip(l, r)
  for (let j = l; j < r; j++) {
    nums[j] ^= 1
  }
  for (let j = 1; j <= n; j++) {
    if (seg.onesCount(j - 1, j) !== nums[j - 1]) {
      console.log('error')
      console.log(j)
      console.log(seg.onesCount(j - 1, j))
      console.log(nums[j - 1])
      console.log(seg.toString())
      throw new Error()
    }
  }
}
