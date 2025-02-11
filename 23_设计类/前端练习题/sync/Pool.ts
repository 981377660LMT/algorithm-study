// Go Pool 核心特性:
//   复用对象以减少内存分配和回收
//   Get() 方法获取对象
//   Put() 方法归还对象
//   New() 函数创建新对象
//
// 实现思路:
//   使用数组作为对象池
//   提供对象获取和归还接口
//   支持自定义对象创建函数
//   保证线程安全

import { Mutex } from './Mutex'

export class Pool<T> {
  private readonly _items: T[] = []
  private readonly _mutex = new Mutex()
  private readonly _newFn: () => T

  constructor(newFn: () => T) {
    this._newFn = newFn
  }

  async get(): Promise<T> {
    await this._mutex.lock()
    try {
      if (this._items.length > 0) {
        return this._items.pop()!
      }
      return this._newFn()
    } finally {
      this._mutex.unlock()
    }
  }

  async put(item: T): Promise<void> {
    await this._mutex.lock()
    try {
      this._items.push(item)
    } finally {
      this._mutex.unlock()
    }
  }

  async size(): Promise<number> {
    await this._mutex.lock()
    try {
      return this._items.length
    } finally {
      this._mutex.unlock()
    }
  }
}

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  async function example() {
    // 创建一个字符串缓冲区对象池
    const pool = new Pool<Buffer>(() => Buffer.alloc(1024))

    // 使用对象
    const buf = await pool.get()
    try {
      // 使用buf
      buf.write('hello')
    } finally {
      // 归还对象
      await pool.put(buf)
    }
  }

  example()
}
