import { SortedList } from '../SortedList'

describe('SortedList', () => {
  let sl: SortedList<any>
  let sortedNums: number[]
  beforeEach(() => {
    const n = Math.floor(Math.random() * 100) + 1
    const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 100000))
    sl = new SortedList(nums)
    sortedNums = nums.sort((a, b) => a - b)
  })

  it('should support add and discard', () => {
    const add = Math.floor(Math.random() * 1000)
    sl.add(add)
    sortedNums.push(add)
    sortedNums.sort((a, b) => a - b)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])

    const discard = Math.floor(Math.random() * 1000)
    sl.discard(discard)
    sortedNums = sortedNums.filter(num => num !== discard)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])
  })

  it('should support bisectLeft and bisectRight', () => {
    const target = Math.floor(Math.random() * 1000)
    const left = sl.bisectLeft(target)
    const right = sl.bisectRight(target)
    expect(left).toBe(sortedNums.findIndex(num => num >= target))
    expect(right).toBe(sortedNums.findIndex(num => num > target))
  })

  it('should support at', () => {
    const index = Math.floor(Math.random() * sl.length)
    expect(sl.at(index)).toBe(sortedNums[index])
  })

  it('should support size', () => {
    expect(sl.length).toBe(sortedNums.length)
  })

  it('should support has', () => {
    const target = Math.floor(Math.random() * 1000)
    expect(sl.has(target)).toBe(sortedNums.includes(target))
  })

  it('should support pop', () => {
    const pop = sl.pop()
    expect(pop).toBe(sortedNums.pop())
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])
  })

  it('should support clear', () => {
    sl.clear()
    expect(sl.length).toBe(0)
    expect([...sl]).toStrictEqual([])
  })

  it('should support entries, forEach, [Symbol.iterator]', () => {
    sl.forEach((value, index) => {
      expect(value).toBe(sortedNums[index])
    })
    for (const [index, item] of sl.entries()) {
      expect(item).toBe(sortedNums[index])
    }
    let pos = 0
    for (const item of sl) {
      expect(item).toBe(sortedNums[pos++])
    }
  })
})
