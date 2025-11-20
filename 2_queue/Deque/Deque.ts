/**
 * 双栈实现的可索引 deque (均摊 O(1)).
 */
export class Deque<T> {
  private _stack1: T[] = []
  private _stack2: T[] = []

  constructor(iterable?: Iterable<T>) {
    if (iterable) {
      for (const x of iterable) this._stack2.push(x)
    }
  }

  push(x: T): void {
    this._stack2.push(x)
  }

  unshift(x: T): void {
    this._stack1.push(x)
  }

  pop(): T | undefined {
    if (this.isEmpty()) return undefined
    if (!this._stack2.length) this._rebalanceFromLeft()
    return this._stack2.pop()!
  }

  shift(): T | undefined {
    if (this.isEmpty()) return undefined
    if (!this._stack1.length) this._rebalanceFromRight()
    return this._stack1.pop()!
  }

  front(): T | undefined {
    return this.isEmpty() ? undefined : this._stack1[this._stack1.length - 1]
  }

  back(): T | undefined {
    return this.isEmpty() ? undefined : this._stack2[this._stack2.length - 1]
  }

  clear(): void {
    this._stack1.length = 0
    this._stack2.length = 0
  }

  reverse(): void {
    const tmp = this._stack1
    this._stack1 = this._stack2
    this._stack2 = tmp
  }

  get(i: number): T | undefined {
    if (i < 0) i += this.length
    if (i < 0 || i >= this.length) return undefined
    const len1 = this._stack1.length
    if (i < len1) {
      return this._stack1[len1 - 1 - i]
    }
    return this._stack2[i - len1]
  }

  set(i: number, value: T): void {
    if (i < 0) i += this.length
    if (i < 0 || i >= this.length) return
    const len1 = this._stack1.length
    if (i < len1) {
      this._stack1[len1 - 1 - i] = value
    } else {
      this._stack2[i - len1] = value
    }
  }

  isEmpty(): boolean {
    return !this._stack1.length && !this._stack2.length
  }

  toArray(): T[] {
    return [...this]
  }

  forEach(callback: (value: T, index: number) => void): void {
    let ptr = 0
    for (let i = this._stack1.length - 1; ~i; --i) {
      callback(this._stack1[i], ptr++)
    }
    for (let i = 0; i < this._stack2.length; ++i) {
      callback(this._stack2[i], ptr++)
    }
  }

  get length(): number {
    return this._stack1.length + this._stack2.length
  }

  *[Symbol.iterator](): Iterator<T> {
    for (let i = this._stack1.length - 1; ~i; --i) {
      yield this._stack1[i]
    }
    for (let i = 0; i < this._stack2.length; ++i) {
      yield this._stack2[i]
    }
  }

  private _rebalanceFromRight(): void {
    const m = this._stack2.length
    if (!m) return
    const k = (m + 1) >>> 1
    this._stack1 = this._stack2.splice(0, k).reverse()
  }

  private _rebalanceFromLeft(): void {
    const m = this._stack1.length
    if (!m) return
    const k = (m + 1) >>> 1
    this._stack2 = this._stack1.splice(0, k).reverse()
  }
}

/**
 * 随机对拍测试
 */
function runTests(): void {
  /**
   * 参考实现：基于单数组的简单 deque，用于对拍。
   */
  class RefDeque<T> {
    private arr: T[] = []
    append(x: T) {
      this.arr.push(x)
    }
    appendleft(x: T) {
      this.arr.unshift(x)
    }
    pop(): T {
      return this.arr.pop() as T
    }
    popleft(): T {
      return this.arr.shift() as T
    }
    get(i: number): T {
      const n = this.arr.length
      if (i < 0) i += n
      if (i < 0 || i >= n) throw new Error('index out of range')
      return this.arr[i]
    }
    set(i: number, v: T) {
      const n = this.arr.length
      if (i < 0) i += n
      if (i < 0 || i >= n) throw new Error('index out of range')
      this.arr[i] = v
    }
    reverse() {
      this.arr.reverse()
    }
    clear() {
      this.arr.length = 0
    }
    toArray() {
      return this.arr.slice()
    }
    get length() {
      return this.arr.length
    }
  }

  const my = new Deque<number>()
  const ref = new RefDeque<number>()

  const NUM_OPERATIONS = 2e6 // 可调整
  const MAX_VALUE = 10_000

  function assert(cond: boolean, msg: string) {
    if (!cond) throw new Error(msg)
  }

  console.log(`开始进行 ${NUM_OPERATIONS} 次随机操作对拍测试...`)

  for (let i = 0; i < NUM_OPERATIONS; i++) {
    const r = Math.random()

    function check(name: string) {
      assert(my.length === ref.length, `操作 ${name} 长度不一致`)
      const a1 = my.toArray()
      const a2 = ref.toArray()
      assert(
        a1.length === a2.length && a1.every((v, idx) => v === a2[idx]),
        `操作 ${name} 内容不一致`
      )
    }

    if (r < 0.25) {
      // append
      const v = Math.floor(Math.random() * MAX_VALUE)
      my.push(v)
      ref.append(v)
      check('append')
    } else if (r < 0.5) {
      // appendleft
      const v = Math.floor(Math.random() * MAX_VALUE)
      my.unshift(v)
      ref.appendleft(v)
      check('appendleft')
    } else if (r < 0.6) {
      // pop
      if (ref.length > 0) {
        const a = my.pop()
        const b = ref.pop()
        assert(a === b, 'pop 返回值不一致')
        check('pop')
      } else {
        let threw = false
        try {
          my.pop()
        } catch {
          threw = true
        }
        // assert(threw, '空队列 pop 未抛出')
      }
    } else if (r < 0.7) {
      // popleft
      if (ref.length > 0) {
        const a = my.shift()
        const b = ref.popleft()
        assert(a === b, 'popleft 返回值不一致')
        check('popleft')
      } else {
        let threw = false
        try {
          my.shift()
        } catch {
          threw = true
        }
        // assert(threw, '空队列 popleft 未抛出')
      }
    } else if (r < 0.8) {
      // get/set
      if (ref.length > 0) {
        let idx = Math.floor(Math.random() * ref.length)
        if (Math.random() < 0.5) idx = idx - ref.length // 负索引测试
        assert(my.get(idx) === ref.get(idx), `get 索引 ${idx} 不一致`)
        if (Math.random() < 0.3) {
          const v = Math.floor(Math.random() * MAX_VALUE)
          my.set(idx, v)
          ref.set(idx, v)
          check(`set at ${idx}`)
        }
      }
    } else if (r < 0.85) {
      // reverse
      if (ref.length > 0) {
        my.reverse()
        ref.reverse()
        check('reverse')
      }
    } else if (r < 0.9) {
      // clear
      my.clear()
      ref.clear()
      check('clear')
    } else {
      // 混合操作触发 rebalance
      const steps = 5 + Math.floor(Math.random() * 16)
      for (let t = 0; t < steps; t++) {
        if (Math.random() < 0.5) {
          const v = Math.floor(Math.random() * MAX_VALUE)
          my.push(v)
          ref.append(v)
        } else {
          if (ref.length > 0) {
            my.shift()
            ref.popleft()
          }
        }
      }
      check('mixed')
    }

    if ((i + 1) % (NUM_OPERATIONS / 10) === 0) {
      console.log(`  已完成 ${i + 1} / ${NUM_OPERATIONS} 次操作...`)
    }
  }

  console.log('所有测试通过！MyDeque 与参考实现行为一致。')
}

// 直接运行
if (require.main === module) {
  runTests()
}
