import { SortedList } from '../SortedList'

export {}

for (let i = 0; i < 100; i++) {
  const n = ~~(Math.random() * 10000) + 1000
  const nums = Array(n)
    .fill(0)
    .map(() => ~~(Math.random() * 100000))
  const sl = new SortedList(nums)
  nums.sort((a, b) => a - b)

  let start = ~~(Math.random() * sl.length)
  let end = ~~(Math.random() * sl.length)
  if (start > end) {
    start ^= end
    end ^= start
    start ^= end
  }
  sl.erase(start, end)
  nums.splice(start, end - start)
  if (nums.toString() !== [...sl].toString()) {
    console.log(nums.toString())
    console.log([...sl].toString())
    console.log(nums.toString() === [...sl].toString())
  }
}
console.log('ok')
