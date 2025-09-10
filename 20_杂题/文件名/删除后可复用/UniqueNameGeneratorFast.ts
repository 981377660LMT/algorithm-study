type NamingStrategy = (baseName: string, suffix: number) => string

const defaultNamingStrategy: NamingStrategy = (baseName, suffix) => `${baseName}(${suffix})`

interface IUniqueNameGeneratorOptions {
  initialNames?: string[]
  namingStrategy?: NamingStrategy
}

export class UniqueNameGeneratorFast {
  private readonly _existingNames = new Set<string>()
  private readonly _pools = new Map<string, IDPool>()
  private readonly _namingStrategy: NamingStrategy

  constructor(options: IUniqueNameGeneratorOptions = {}) {
    const { initialNames = [], namingStrategy = defaultNamingStrategy } = options
    this._namingStrategy = namingStrategy
    initialNames.forEach(name => this.add(name))
  }

  add(name: string): string {
    if (!this._pools.has(name)) this._pools.set(name, new IDPool())
    const suffix = this._pools.get(name)!.alloc()
    const newName = suffix ? this._namingStrategy(name, suffix) : name
    this._existingNames.add(newName)
    return newName
  }

  delete(name: string, suffix: number): boolean {
    if (!this._existingNames.has(name)) {
      return false
    }
    this._existingNames.delete(name)
    const pool = this._pools.get(name)!
    pool.release(suffix)
    if (!pool.size) this._pools.delete(name)
    return true
  }

  has(name: string): boolean {
    return this._existingNames.has(name)
  }

  getAll(): string[] {
    return Array.from(this._existingNames)
  }
}

class IDPool {
  private readonly _reused = new Heap<number>([], (a, b) => a < b)
  private _nextId = 0

  constructor(startId = 0) {
    this._nextId = startId
  }

  alloc(): number {
    if (this._reused.length) {
      return this._reused.pop()!
    }
    return this._nextId++
  }

  release(id: number): void {
    this._reused.push(id)
  }

  reset(): void {
    this._reused.clear()
    this._nextId = 0
  }

  get size(): number {
    return this._nextId - this._reused.length
  }
}

class Heap<T = any> {
  private readonly _data: T[]
  private readonly _less: (a: T, b: T) => boolean

  constructor(data: T[], less: (a: T, b: T) => boolean) {
    this._data = data
    this._less = less
    if (data.length > 1) {
      this._heapify()
    }
  }

  peek(): T | undefined {
    return this._data[0]
  }

  push(x: T): void {
    this._data.push(x)
    this._up(this.length - 1)
  }

  pop(): T | undefined {
    if (!this.length) {
      return undefined
    }
    this._swap(0, this.length - 1)
    const res = this._data.pop()
    if (this.length) {
      this._down(0)
    }
    return res
  }

  clear(): void {
    this._data.length = 0
  }

  get length(): number {
    return this._data.length
  }

  private _heapify(): void {
    const n = this.length
    for (let i = (n >> 1) - 1; ~i; i--) {
      this._down(i)
    }
  }

  private _swap(i: number, j: number): void {
    const tmp = this._data[i]
    this._data[i] = this._data[j]
    this._data[j] = tmp
  }

  private _up(j: number): void {
    const data = this._data
    const less = this._less
    const item = data[j]
    while (j > 0) {
      const i = (j - 1) >>> 1 // parent
      if (!less(item, data[i])) {
        break
      }
      data[j] = data[i]
      j = i
    }
    data[j] = item
  }

  private _down(i0: number): void {
    const data = this._data
    const less = this._less
    const n = this.length
    const item = data[i0]
    let i = i0
    while (true) {
      const j1 = (i << 1) | 1
      if (j1 >= n) {
        break
      }
      let j = j1
      const j2 = j1 + 1
      if (j2 < n && less(data[j2], data[j1])) {
        j = j2
      }
      if (!less(data[j], item)) {
        break
      }
      data[i] = data[j]
      i = j
    }
    data[i] = item
  }
}
