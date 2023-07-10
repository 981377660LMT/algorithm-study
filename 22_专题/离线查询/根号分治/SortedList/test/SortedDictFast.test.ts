import { SortedDictFast } from '../SortedDictFast'

describe('SortedDictFast', () => {
  let sd: SortedDictFast<number, number>
  let sortedDict: Map<number, number>
  beforeEach(() => {
    const n = Math.floor(Math.random() * 100) + 1
    const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 100000))
    nums.sort((a, b) => a - b)
    sd = new SortedDictFast(nums.map(num => [num, num]))
    sortedDict = new Map(nums.map(num => [num, num]))
  })

  // set/setDefault
  it('should support set and setDefault', () => {
    const key = Math.floor(Math.random() * 1000)
    const value = Math.floor(Math.random() * 1000)
    sd.set(key, value)
    sortedDict.set(key, value)
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))

    const key2 = Math.floor(Math.random() * 1000)
    const value2 = Math.floor(Math.random() * 1000)
    sd.setDefault(key2, value2)
    if (!sortedDict.has(key2)) sortedDict.set(key2, value2)
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))
  })

  // has
  it('should support has', () => {
    const key = Math.floor(Math.random() * 1000)
    expect(sd.has(key)).toBe(sortedDict.has(key))
  })

  // get
  it('should support get', () => {
    const key = Math.floor(Math.random() * 1000)
    expect(sd.get(key)).toBe(sortedDict.get(key))
  })

  // delete/pop
  it('should support delete and pop', () => {
    const key = Math.floor(Math.random() * 1000)
    sd.delete(key)
    sortedDict.delete(key)
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))

    const key2 = Math.floor(Math.random() * 1000)
    const popped = sd.pop(key2)
    const sortedPop = sortedDict.get(key2)
    sortedDict.delete(key2)
    expect(popped).toBe(sortedPop)
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))
  })

  // popItem
  it('should support popItem', () => {
    const popped = sd.popItem(0)
    const sortedPop = sortedDict.entries().next().value
    sortedDict.delete(sortedPop[0])
    expect(popped).toStrictEqual(sortedPop)
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))
  })

  // peekItem/peekMinItem/peekMaxItem
  it('should support peekItem', () => {
    const peeked = sd.peekItem(0)
    const sortedPeek = sortedDict.entries().next().value
    expect(peeked).toStrictEqual(sortedPeek)

    const peekedMin = sd.peekMinItem()
    const sortedPeekMin = [...sortedDict].sort((a, b) => a[0] - b[0])[0]
    expect(peekedMin).toStrictEqual(sortedPeekMin)

    const peekedMax = sd.peekMaxItem()
    const sortedPeekMax = [...sortedDict].sort((a, b) => b[0] - a[0])[0]
    expect(peekedMax).toStrictEqual(sortedPeekMax)
  })

  // forEach
  it('should support forEach', () => {
    const callback = jest.fn()
    sd.forEach(callback)
    expect(callback).toBeCalledTimes(sortedDict.size)
    expect(callback).toBeCalledWith(...sortedDict.entries().next().value)

    const items1: [number, number][] = []
    const items2: [number, number][] = []
    sd.forEach((value, key) => items1.push([key, value]))
    sortedDict.forEach((value, key) => items2.push([key, value]))
    expect(items1).toStrictEqual(items2)
  })

  // enumerate
  it('should support enumerate', () => {
    const items1: [number, number][] = []
    const end = Math.floor(sd.size * Math.random())
    sd.enumerate(0, end, (value, key) => items1.push([key, value]))
    const items2 = [...sortedDict].slice(0, end).sort((a, b) => a[0] - b[0])
    expect(items1).toStrictEqual(items2)
  })

  // bisectLeft/bisectRight
  it('should support bisectLeft and bisectRight', () => {
    const key = Math.floor(Math.random() * 1000)
    const left = sd.bisectLeft(key)
    const right = sd.bisectRight(key)
    const sortedLeft = [...sortedDict.keys()].filter(k => k < key).length
    const sortedRight = [...sortedDict.keys()].filter(k => k <= key).length
    expect(left).toBe(sortedLeft)
    expect(right).toBe(sortedRight)
  })

  // floor/ceiling/lower/higher
  it('should support floor, ceiling, lower and higher', () => {
    const key = Math.floor(Math.random() * 1000)
    const floor = sd.floor(key)
    const ceiling = sd.ceiling(key)
    const lower = sd.lower(key)
    const higher = sd.higher(key)
    const sortedFloor = [...sortedDict.entries()].filter(item => item[0] <= key).pop()
    const sortedCeiling = [...sortedDict.entries()].filter(item => item[0] >= key).shift()
    const sortedLower = [...sortedDict.entries()].filter(item => item[0] < key).pop()
    const sortedUpper = [...sortedDict.entries()].filter(item => item[0] > key).shift()
    expect(floor).toStrictEqual(sortedFloor)
    expect(ceiling).toStrictEqual(sortedCeiling)
    expect(lower).toStrictEqual(sortedLower)
    expect(higher).toStrictEqual(sortedUpper)
  })

  // clear
  it('should support clear', () => {
    sd.clear()
    sortedDict.clear()
    expect(sd.size).toBe(sortedDict.size)
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))
  })

  // islice
  it('should support islice', () => {
    const start = Math.floor(sd.size * Math.random())
    const end = Math.floor(sd.size * Math.random())
    const reverse = Math.random() > 0.5
    const sliced = sd.islice(start, end, reverse)
    let sortedSliced = [...sortedDict].slice(start, end)
    if (reverse) sortedSliced.reverse()
    expect([...sliced]).toStrictEqual(sortedSliced)
  })

  // irange
  it('should support irange', () => {
    for (let i = 0; i < 10; i++) {
      const min = Math.floor(sd.size * Math.random())
      const max = Math.floor(sd.size * Math.random())
      const reverse = Math.random() > 0.5
      const ranged = sd.irange(min, max, reverse)
      const sortedRanged = [...sortedDict]
        .sort((a, b) => a[0] - b[0])
        .filter(([k]) => k >= min && k <= max)
      if (reverse) sortedRanged.reverse()
      expect([...ranged]).toStrictEqual(sortedRanged)
    }
  })

  // keys/values/entries/[Symbol.iterator]
  it('should support keys, values, entries and [Symbol.iterator]', () => {
    expect([...sd.keys()]).toStrictEqual([...sortedDict.keys()])
    expect([...sd.values()]).toStrictEqual([...sortedDict.values()])
    expect([...sd.entries()]).toStrictEqual([...sortedDict.entries()])
    expect([...sd]).toStrictEqual([...sortedDict].sort((a, b) => a[0] - b[0]))
  })
})
