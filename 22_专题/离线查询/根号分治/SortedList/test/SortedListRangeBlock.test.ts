import { SortedListRangeBlock } from '../../值域分块/SortedListRangeBlock'

describe('SortedListFast', () => {
  let sl: SortedListRangeBlock
  let sortedNums: number[]
  beforeEach(() => {
    const n = Math.floor(Math.random() * 100) + 1
    const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 1e5))
    sl = new SortedListRangeBlock(1e6)
    sl.update(...nums)
    sortedNums = nums.sort((a, b) => a - b)
  })

  it('should support add and discard', () => {
    const add = Math.floor(Math.random() * 100)
    sl.add(add)
    sortedNums.push(add)
    sortedNums.sort((a, b) => a - b)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])

    const discard = Math.floor(Math.random() * 100)
    sl.discard(discard)
    sortedNums = sortedNums.filter(num => num !== discard)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])
  })

  it('should support bisectLeft and bisectRight', () => {
    const target = Math.floor(Math.random() * 100)
    const left = sl.bisectLeft(target)
    const right = sl.bisectRight(target)
    const leftIndex = sortedNums.findIndex(num => num >= target)
    if (leftIndex === -1) {
      expect(left).toBe(sl.length)
    } else {
      expect(left).toBe(leftIndex)
    }
    const rightIndex = sortedNums.findIndex(num => num > target)
    if (rightIndex === -1) {
      expect(right).toBe(sl.length)
    } else {
      expect(right).toBe(rightIndex)
    }
  })

  it('should support at', () => {
    const index = Math.floor(Math.random() * sl.length)
    expect(sl.at(index)).toBe(sortedNums[index])
  })

  it('should support size', () => {
    expect(sl.length).toBe(sortedNums.length)
  })

  it('should support has', () => {
    const target = Math.floor(Math.random() * 100)
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

  it('should support erase', () => {
    for (let i = 0; i < 100; i++) {
      let start = Math.floor(Math.random() * sl.length)
      let end = Math.floor(Math.random() * sl.length)
      sl.erase(start, end)
      sortedNums.splice(start, end - start)
      expect(sl.length).toBe(sortedNums.length)
      expect(sortedNums).toStrictEqual([...sl])
    }
  })

  it('should support slice', () => {
    for (let i = 0; i < 100; i++) {
      let start = Math.floor(Math.random() * sl.length)
      let end = Math.floor(Math.random() * sl.length)
      const reverse = Math.random() > 0.5

      const slice = sl.slice(start, end)
      if (reverse) slice.reverse()

      expect(slice).toStrictEqual(
        reverse ? sortedNums.slice(start, end).reverse() : sortedNums.slice(start, end)
      )
    }
  })

  // enumerate
  it('should support enumerate', () => {
    for (let i = 0; i < 100; i++) {
      let start = Math.floor(Math.random() * sl.length)
      let end = Math.floor(Math.random() * sl.length)
      const enumerated: [number][] = []
      sl.enumerate(start, end, v => {
        enumerated.push([v])
      })
      const target = sortedNums.slice(start, end).map((v, i) => [v])
      expect(enumerated).toStrictEqual(target)
    }
  })

  // update
  it('should support update', () => {
    // update(small)
    const small = Array(10)
      .fill(0)
      .map(() => Math.floor(Math.random() * 100))
    sl.update(...small)
    sortedNums.push(...small)
    sortedNums.sort((a, b) => a - b)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])

    const big = Array(sl.length * 10)
      .fill(0)
      .map(() => Math.floor(Math.random() * 100))
    sl.update(...big)
    sortedNums.push(...big)
    sortedNums.sort((a, b) => a - b)
    expect(sl.length).toBe(sortedNums.length)
    expect(sortedNums).toStrictEqual([...sl])
  })

  // merge
  it('should support merge', () => {
    const sl2 = new SortedListRangeBlock(1e6)
    const nums2 = Array(100)
      .fill(0)
      .map(() => Math.floor(Math.random() * 100))
    sl2.update(...nums2)
    sl.merge(sl2)
    const target = [...sortedNums, ...nums2].sort((a, b) => a - b)
    expect(sl.length).toBe(target.length)
    expect([...sl]).toStrictEqual(target)
  })

  // sumSlice(start, end)
  it('should support sumSlice', () => {
    for (let i = 0; i < 100; i++) {
      let start = Math.floor(Math.random() * sl.length)
      let end = Math.floor(Math.random() * sl.length)
      if (start > end) [start, end] = [end, start]
      const sum = sl.sumSlice(start, end)
      const target = sortedNums.slice(start, end).reduce((a, b) => a + b, 0)
      expect(sum).toBe(target)

      // discard/add
      const willDiscard = Math.random() > 0.5
      if (willDiscard) {
        const discard = sortedNums[~~(Math.random() * sortedNums.length)]
        sl.discard(discard)
        const index = sortedNums.findIndex(num => num === discard)
        sortedNums.splice(index, 1)
      } else {
        const add = Math.floor(Math.random() * 100)
        sl.add(add)
        sortedNums.push(add)
        sortedNums.sort((a, b) => a - b)
      }
    }
  })

  // sumRange(min, max)
  it('should support sumRange', () => {
    for (let i = 0; i < 100; i++) {
      let min = Math.floor(Math.random() * 100000)
      let max = Math.floor(Math.random() * 100000)
      if (min > max) [min, max] = [max, min]
      const sum = sl.sumRange(min, max)
      const target = sortedNums.filter(num => num >= min && num <= max).reduce((a, b) => a + b, 0)
      expect(sum).toBe(target)

      // discard/add
      const willDiscard = Math.random() > 0.5
      if (willDiscard) {
        const randint = sortedNums[~~(Math.random() * sortedNums.length)]
        const index = sortedNums.findIndex(num => num === randint)
        sl.discard(randint)
        sortedNums.splice(index, 1)
      } else {
        const randint = Math.floor(Math.random() * 100)
        sl.add(randint)
        sortedNums.push(randint)
        sortedNums.sort((a, b) => a - b)
      }
    }
  })
})
