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

// irange
for (let i = 0; i < 100; i++) {
  const n = ~~(Math.random() * 10000) + 1000
  const nums = Array(n)
    .fill(0)
    .map(() => ~~(Math.random() * 100000))
  const sl = new SortedList(nums)
  nums.sort((a, b) => a - b)
  let min = Math.floor(Math.random() * sl.length)
  let max = Math.floor(Math.random() * sl.length)
  const irange = [...sl.irange(min, max)]
  const target = nums.filter(num => num >= min && num <= max).sort((a, b) => a - b)
  if (irange.toString() !== target.toString()) {
    console.log(irange.toString())
    console.log(target.toString())
    console.log(irange.toString() === target.toString())
  }
}

// islice
for (let i = 0; i < 100; i++) {
  const n = ~~(Math.random() * 1) + 10
  const nums = Array(n)
    .fill(0)
    .map(() => ~~(Math.random() * 100000))
  const sl = new SortedList(nums)
  nums.sort((a, b) => a - b)

  let start = Math.floor(Math.random() * sl.length)
  let end = Math.floor(Math.random() * sl.length)
  const reverse = Math.random() > 0.5
  const islice = [...sl.islice(start, end, reverse)]
  const slice = sl.slice(start, end)
  if (reverse) slice.reverse()

  if (islice.toString() !== slice.toString()) {
    console.log(islice.toString())
    console.log(slice.toString())
    console.log(islice.toString() === slice.toString())
    console.log(sl.toString(), start, end, reverse)
    break
  }
}

console.log('ok')
