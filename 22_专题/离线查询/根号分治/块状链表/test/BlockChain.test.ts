import { AbstractBlock, BlockChain } from '../BlockChain'

class NumberBlockArray extends AbstractBlock<{ value: number }, number, number, NumberBlockArray> {
  private _data: number[]
  private _sum: number
  private _lazy = 0

  constructor(data?: number[]) {
    super()
    this._data = data ? data : []
    let sum = 0
    for (let i = 0; i < this._data.length; i++) sum += this._data[i]
    this._sum = sum
  }

  override split(k: number): { block1: NumberBlockArray; block2: NumberBlockArray } {
    const remain = this._data.splice(k, this._data.length - k)
    const block1 = new NumberBlockArray(this._data)
    block1._lazy = this._lazy
    const block2 = new NumberBlockArray(remain)
    block2._lazy = this._lazy
    return { block1, block2 }
  }

  override merge(other: NumberBlockArray): NumberBlockArray {
    // this._data.push(...other._data)
    // const lazyDiff = other._lazy - this._lazy
    // for (let i = 0; i < other._data.length; i++) this._sum += other._data[i] + lazyDiff
    // return this
    const newData = Array(this._data.length + other._data.length)
    for (let i = 0; i < this._data.length; i++) newData[i] = this._data[i] + this._lazy
    for (let i = 0; i < other._data.length; i++)
      newData[this._data.length + i] = other._data[i] + other._lazy
    return new NumberBlockArray(newData)
  }

  override insertBefore(index: number, e: number): void {
    this._data.splice(index, 0, e)
    this._pushDown()
  }

  override delete(index: number): void {
    this._data.splice(index, 1)
    // this._sum -= this._lazy
    this._pushDown()
  }

  override get(index: number): number {
    return this._data[index] + this._lazy
  }

  override reverse(): void {
    this._data.reverse()
  }

  override fullyQuery(sum: { value: number }): void {
    sum.value += this._sum + this._lazy * this._data.length
  }

  override partialQuery(index: number, sum: { value: number }): void {
    sum.value += this._data[index] + this._lazy
  }

  override fullyUpdate(lazy: number): void {
    this._lazy += lazy
  }

  override partialUpdate(index: number, lazy: number): void {
    this._data[index] += lazy
    this._sum += lazy
  }

  override beforePartialQuery(): void {
    this._pushDown()
  }

  override afterPartialUpdate(): void {
    this._pushDown()
  }

  private _pushDown(): void {
    this._sum = 0
    for (let i = 0; i < this._data.length; i++) {
      this._data[i] += this._lazy
      this._sum += this._data[i]
    }
    this._lazy = 0
  }
}

// get/prefixSize/split/mergeDestructively/insetBefore/delete
// reverse/rotateLeft/getAll/length/query/update
describe('BlockChain', () => {
  let nums: number[]
  let chain: BlockChain<{ value: number }, number, number, NumberBlockArray>

  beforeEach(() => {
    // nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    nums = Array.from({ length: 100 }, () => ~~(Math.random() * 100))
    chain = BlockChain.create(
      nums.length,
      (start, end) => new NumberBlockArray(nums.slice(start, end))
    )
  })

  test('get', () => {
    for (let i = 0; i < chain.length; i++) expect(chain.get(i)).toBe(nums[i])
  })

  test('prefixSize', () => {
    let size = 0
    chain.enumerateNode(node => {
      expect(chain.prefixSize(node.data!, false)).toBe(size)
      size += node.size
      expect(chain.prefixSize(node.data!, true)).toBe(size)
    })
  })

  test('split', () => {
    const index = ~~(Math.random() * (nums.length + 1))
    const { first, second } = chain.split(index, () => new NumberBlockArray())
    expect(first.getAll()).toEqual(nums.slice(0, index))
    expect(second.getAll()).toEqual(nums.slice(index))
  })

  test('mergeDestructively', () => {
    const index = ~~(Math.random() * (nums.length + 1))
    const { first, second } = chain.split(index, () => new NumberBlockArray())
    const firstData = first.getAll()
    const secondData = second.getAll()
    first.mergeDestructively(second)
    expect(first.getAll()).toEqual([...firstData, ...secondData])
  })

  test('insertBefore/delete', () => {
    for (let _ = 0; _ < 100; _++) {
      const index = ~~(Math.random() * (nums.length + 1))
      const num = ~~(Math.random() * 100)
      nums.splice(index, 0, num)
      chain.insertBefore(index, num)
      expect(chain.getAll()).toEqual(nums)
      if (Math.random() < 0.5) {
        const index = ~~(Math.random() * nums.length)
        nums.splice(index, 1)
        chain.delete(index)
        expect(chain.getAll()).toEqual(nums)
      }
    }
  })

  test('reverse', () => {
    for (let _ = 0; _ < 100; _++) {
      const start = ~~(Math.random() * nums.length)
      const end = start + ~~(Math.random() * (nums.length - start + 1))
      chain.reverse(start, end)
      for (let i = start, j = end - 1; i < j; i++, j--) {
        const temp = nums[i]
        nums[i] = nums[j]
        nums[j] = temp
      }
      expect(chain.getAll()).toEqual(nums)
    }
  })

  test('rotateLeft', () => {
    for (let _ = 0; _ < 100; _++) {
      const step = ~~(Math.random() * nums.length * 10)
      chain.rotateLeft(step)
      const tmp = ((step % nums.length) + nums.length) % nums.length
      nums = nums.slice(tmp).concat(nums.slice(0, tmp))
      expect(chain.getAll()).toEqual(nums)
    }
  })

  test('query/update', () => {
    for (let _ = 0; _ < 100; _++) {
      const start = ~~(Math.random() * nums.length)
      const end = start + ~~(Math.random() * (nums.length - start + 1))
      const sum = { value: 0 }
      chain.query(start, end, sum)
      expect(sum.value).toBe(nums.slice(start, end).reduce((acc, cur) => acc + cur, 0))
      const lazy = ~~(Math.random() * 100)
      chain.update(start, end, lazy)
      for (let i = start; i < end; i++) nums[i] += lazy
      expect(chain.getAll()).toEqual(nums)
    }
  })

  test('query/update with insetBefore/delete/reverse/rotateLeft', () => {
    for (let _ = 0; _ < 2000; _++) {
      const start = ~~(Math.random() * nums.length)
      const end = start + ~~(Math.random() * (nums.length - start + 1))
      const sum = { value: 0 }
      chain.query(start, end, sum)
      expect(sum.value).toBe(nums.slice(start, end).reduce((acc, cur) => acc + cur, 0))

      // const lazy = ~~(Math.random() * 100)
      // chain.update(start, end, lazy)
      // for (let i = start; i < end; i++) nums[i] += lazy
      // expect(chain.getAll()).toEqual(nums)

      // insert
      const insertIndex = ~~(Math.random() * (nums.length + 1))
      const insertNum = ~~(Math.random() * 100)
      nums.splice(insertIndex, 0, insertNum)
      chain.insertBefore(insertIndex, insertNum)

      // delete
      const deleteIndex = ~~(Math.random() * nums.length)
      nums.splice(deleteIndex, 1)
      chain.delete(deleteIndex)

      // reverse
      const reverseStart = ~~(Math.random() * nums.length)
      const reverseEnd = reverseStart + ~~(Math.random() * (nums.length - reverseStart + 1))
      chain.reverse(reverseStart, reverseEnd)
      nums = nums
        .slice(0, reverseStart)
        .concat(nums.slice(reverseStart, reverseEnd).reverse(), nums.slice(reverseEnd))

      // rotateLeft
      const step = ~~(Math.random() * nums.length * 10)
      chain.rotateLeft(step)
      const tmp = ((step % nums.length) + nums.length) % nums.length
      nums = nums.slice(tmp).concat(nums.slice(0, tmp))
    }
  })
})
