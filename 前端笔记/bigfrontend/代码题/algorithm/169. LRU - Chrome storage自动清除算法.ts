/* eslint-disable prefer-destructuring */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 169. LRU - Chrome storage自动清除算法

interface OriginData {
  origin: string
  lastUsed: number
  size: number
  persistent: boolean
}

interface LRUStorage {
  capacity: number

  // !getData 和setData被调用的时候，需要算作'used'。
  // to use the data for origin
  // return size of the data or undefined if not exist
  getData(origin: string): OriginData | undefined

  // updating data for origin
  // return boolean to indicate success or failure
  // If the total size exceeds capacity,
  // !Least Recently Used non-persistent origin data other than itself should be evicted.
  setData(origin: string, size: number): boolean

  // manually clear data for origin
  clearData(origin: string): void

  // change data for origin to be persistent
  // it only handles existing data not the data added later
  // !persistent data cannot be evicted unless manually clear it
  makePersistent(origin: string): void
}

class MyLRUStorage implements LRUStorage {
  readonly capacity: number
  private _normalSize = 0
  private _persistedSize = 0
  private readonly _getTimestamp: () => number
  private readonly _normalData = new Map<string, OriginData>()
  private readonly _persistedData = new Map<string, OriginData>()

  // !由于时间精度问题，class的constructor请支持第二个getTimestamp 参数。
  constructor(capacity: number, getTimestamp: () => number) {
    this.capacity = capacity
    this._getTimestamp = getTimestamp
  }

  getData(origin: string): OriginData | undefined {
    if (!this._normalData.has(origin) && !this._persistedData.has(origin)) return undefined
    if (this._persistedData.has(origin)) return this._persistedData.get(origin)
    const data = this._normalData.get(origin)!
    this._normalData.delete(origin)
    const newData = { ...data, lastUsed: this._getTimestamp() }
    this._normalData.set(origin, newData)
    return newData
  }

  setData(origin: string, size: number): boolean {
    if (size > this.capacity) return false
    if (this._persistedSize + size > this.capacity) return false
    if (this._normalData.has(origin)) {
      this._normalSize -= this._normalData.get(origin)!.size
      this._normalData.delete(origin)
    }

    this._normalSize += size
    this._normalData.set(origin, {
      origin,
      size,
      lastUsed: this._getTimestamp(),
      persistent: false
    })

    // iterate keys to remove the least used data
    if (this.size <= this.capacity) return true
    const keys = this._normalData.keys()
    while (this.size > this.capacity) {
      const leastUsed = keys.next().value
      this._normalSize -= this._normalData.get(leastUsed)!.size
      this._normalData.delete(leastUsed)
    }
    return true
  }

  clearData(origin: string): void {
    if (this._normalData.has(origin)) {
      this._normalSize -= this._normalData.get(origin)!.size
      this._normalData.delete(origin)
    }
    if (this._persistedData.has(origin)) {
      this._persistedSize -= this._persistedData.get(origin)!.size
      this._persistedData.delete(origin)
    }
  }

  makePersistent(origin: string): void {
    if (this._persistedData.has(origin) || !this._normalData.has(origin)) return
    const size = this._normalData.get(origin)!.size
    this._persistedSize += size
    this._normalSize -= size
    this._persistedData.set(origin, this._normalData.get(origin)!)
    this._normalData.delete(origin)
  }

  get size(): number {
    return this._normalSize + this._persistedSize
  }
}

if (require.main === module) {
  const storage = new MyLRUStorage(10, () => 0)
  storage.setData('a', 1)
  storage.setData('b', 3)
  console.dir(storage, { depth: null })
  storage.makePersistent('a')
  storage.makePersistent('b')
  console.dir(storage, { depth: null })

  const result = storage.setData('c', 7)
  console.log(result) // false
  storage.clearData('b')
  storage.setData('c', 7)
  console.log(storage.getData('a')?.size) // 1
  console.log(storage.getData('b')) // undefined
  console.log(storage.getData('c')?.size) // 7
}

export {}
