import { BitVector } from '../BitVector'

// countPrefix
for (let _ = 0; _ < 1000; _++) {
  const n = Math.floor(Math.random() * 100) + 100
  const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 2))
  const bitVector = new BitVector(n)
  for (let i = 0; i < n; i++) {
    if (nums[i]) {
      bitVector.add(i)
    }
  }
  bitVector.build()

  if (
    bitVector.countPrefix(0, bitVector.length) !==
    nums.slice(0, bitVector.length).filter(v => v === 0).length
  ) {
    console.log('error')
    console.log(n, nums)
    console.log(bitVector.countPrefix(0, bitVector.length), bitVector.length)
    console.log(nums.slice(0, bitVector.length).filter(v => v === 0).length)
    console.log(bitVector.toString())
    throw new Error()
  }
}
console.log('pass')
