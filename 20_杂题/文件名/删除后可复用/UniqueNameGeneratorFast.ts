type NamingStrategy = (baseName: string, suffix: number) => string
type ParsingStrategy = (name: string) => [string, number] | null

const defaultNamingStrategy: NamingStrategy = (baseName, suffix) => `${baseName}(${suffix})`

const defaultParsingStrategy: ParsingStrategy = name => {
  const match = name.match(/^(.*)\((\d+)\)$/)
  if (match && match[1] && match[2]) {
    return [match[1], parseInt(match[2], 10)]
  }
  return null
}

interface IUniqueNameGeneratorOptions {
  initialNames?: string[]
  namingStrategy?: NamingStrategy
  parsingStrategy?: ParsingStrategy
}

/**
 * 一个高效的唯一名称生成器，支持删除后复用名称。
 * 内部为每个基础名称维护一个 ID 池（最小堆）来快速获取可用的最小后缀。
 */
export class UniqueNameGeneratorFast {
  private readonly _existingNames = new Set<string>()
  private readonly _pools = new Map<string, IDPool>()
  private readonly _namingStrategy: NamingStrategy
  private readonly _parsingStrategy: ParsingStrategy

  constructor(options: IUniqueNameGeneratorOptions = {}) {
    const {
      initialNames = [],
      namingStrategy = defaultNamingStrategy,
      parsingStrategy = defaultParsingStrategy
    } = options
    this._namingStrategy = namingStrategy
    this._parsingStrategy = parsingStrategy
    initialNames.forEach(name => this.add(name))
  }

  add(name: string): string {
    if (!this._existingNames.has(name)) {
      this._existingNames.add(name)
      return name
    }
    if (!this._pools.has(name)) this._pools.set(name, new IDPool(1))
    const pool = this._pools.get(name)!
    const suffix = pool.alloc()
    const newName = this._namingStrategy(name, suffix)
    this._existingNames.add(newName)
    return newName
  }

  delete(name: string): boolean {
    if (!this._existingNames.delete(name)) {
      return false
    }

    // 尝试解析名称，如果成功，则将后缀释放回对应的池中
    const parsed = this._parsingStrategy(name)
    if (parsed) {
      const [baseName, suffix] = parsed
      const pool = this._pools.get(baseName)
      if (pool) {
        pool.release(suffix)
      }
    }
    return true
  }

  has(name: string): boolean {
    return this._existingNames.has(name)
  }

  rename(oldName: string, newName: string): string | undefined {
    if (!this.has(oldName)) {
      return undefined
    }
    this.delete(oldName)
    return this.add(newName)
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

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  function runTests() {
    console.log('--- 开始测试 ---')

    // 测试 1: 基本的添加和冲突处理
    console.log('\n--- 测试 1: 基本的添加和冲突处理 ---')
    const gen1 = new UniqueNameGeneratorFast()
    console.assert(gen1.add('file') === 'file', 'Test 1.1 failed')
    console.assert(gen1.add('file') === 'file(1)', 'Test 1.2 failed')
    console.assert(gen1.add('image') === 'image', 'Test 1.3 failed')
    console.assert(gen1.add('file') === 'file(2)', 'Test 1.4 failed')
    console.assert(gen1.has('file(1)') === true, 'Test 1.5 failed')
    console.assert(gen1.has('file(3)') === false, 'Test 1.6 failed')
    console.log('当前所有名称:', gen1.getAll())

    // 测试 2: 删除和后缀复用
    console.log('\n--- 测试 2: 删除和后缀复用 ---')
    const gen2 = new UniqueNameGeneratorFast({
      initialNames: ['doc', 'doc(1)', 'doc(2)', 'doc(3)']
    })
    console.assert(gen2.add('doc') === 'doc(4)', 'Test 2.1 failed')
    gen2.delete('doc(2)')
    console.log('删除 "doc(2)"')
    console.assert(gen2.add('doc') === 'doc(2)', 'Test 2.2 failed: 未复用最小后缀')
    gen2.delete('doc(1)')
    console.log('删除 "doc(1)"')
    gen2.delete('doc(3)')
    console.log('删除 "doc(3)"')
    console.assert(gen2.add('doc') === 'doc(1)', 'Test 2.3 failed: 未复用最小后缀')
    console.assert(gen2.add('doc') === 'doc(3)', 'Test 2.4 failed: 未复用次最小后缀')
    console.assert(gen2.add('doc') === 'doc(5)', 'Test 2.5 failed: 未使用新的后缀')
    console.log('当前所有名称:', gen2.getAll())

    // 测试 3: 重命名 (rename)
    console.log('\n--- 测试 3: 重命名 (rename) ---')
    const gen3 = new UniqueNameGeneratorFast({ initialNames: ['a', 'b'] })
    // 场景 3.1: 重命名为不存在的名称
    const rename1 = gen3.rename('a', 'c')
    console.assert(rename1 === 'c', 'Test 3.1 failed')
    console.assert(gen3.has('a') === false && gen3.has('c') === true, 'Test 3.1 state failed')
    console.log('将 "a" 重命名为 "c" ->', rename1)

    // 场景 3.2: 重命名为已存在的名称
    const rename2 = gen3.rename('c', 'b')
    console.assert(rename2 === 'b(1)', 'Test 3.2 failed')
    console.assert(gen3.has('c') === false && gen3.has('b(1)') === true, 'Test 3.2 state failed')
    console.log('将 "c" 重命名为 "b" ->', rename2)

    // 场景 3.3: 重命名一个带后缀的名称
    gen3.rename('b(1)', 'd')
    console.assert(gen3.has('b(1)') === false && gen3.has('d') === true, 'Test 3.3 failed')
    console.log('将 "b(1)" 重命名为 "d"')

    // 场景 3.4: 验证重命名后后缀是否被回收
    const rename4 = gen3.add('b')
    console.assert(rename4 === 'b(1)', 'Test 3.4 failed: 重命名后后缀未回收')
    console.log('再次添加 "b" ->', rename4, '(验证了后缀 "1" 已被回收复用)')
    console.log('当前所有名称:', gen3.getAll())

    // 测试 4: 自定义策略
    console.log('\n--- 测试 4: 自定义策略 ---')
    const gen4 = new UniqueNameGeneratorFast({
      namingStrategy: (base, suffix) => `${base}_${suffix}`,
      parsingStrategy: name => {
        const match = name.match(/^(.*)_(\d+)$/)
        return match ? [match[1], parseInt(match[2])] : null
      }
    })
    gen4.add('item')
    console.assert(gen4.add('item') === 'item_1', 'Test 4.1 failed')
    gen4.add('item') // item_2
    gen4.delete('item_1')
    console.assert(gen4.add('item') === 'item_1', 'Test 4.2 failed: 自定义策略下复用失败')
    console.log('当前所有名称:', gen4.getAll())

    // 测试 5: 边缘情况
    console.log('\n--- 测试 5: 边缘情况 ---')
    const gen5 = new UniqueNameGeneratorFast()
    // 5.1 删除不存在的元素
    console.assert(gen5.delete('non-existent') === false, 'Test 5.1 failed')
    // 5.2 重命名不存在的元素
    console.assert(gen5.rename('non-existent', 'new-name') === undefined, 'Test 5.2 failed')
    // 5.3 添加一个看起来像带后缀的名称
    gen5.add('file(1)')
    console.assert(gen5.has('file(1)'), 'Test 5.3.1 failed')
    // 5.4 此时再添加 file，应该生成 file(1) 吗？不，因为 file(1) 已被占用
    gen5.add('file')
    console.assert(gen5.add('file(1)') === 'file(1)(1)', 'Test 5.4 failed')
    console.log('当前所有名称:', gen5.getAll())

    console.log('\n--- 所有测试执行完毕 ---')
  }

  runTests()
}
