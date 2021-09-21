interface Iterator {
  hasNext: () => boolean
  next: () => number
}

const None = Symbol('None') as any

// 引入缓存
class PeekingIterator {
  private iterator: Iterator
  private cache: number
  constructor(iterator: Iterator) {
    this.iterator = iterator
    this.cache = None
  }

  peek(): number {
    if (this.cache === None) {
      this.cache = this.iterator.next()
    }
    return this.cache!
  }

  next(): number {
    if (this.cache === None) {
      return this.iterator.next()
    }
    const tmp = this.cache
    this.cache = None
    return tmp
  }

  hasNext(): boolean {
    if (this.cache === None) return this.iterator.hasNext()
    return true
  }
}

export {}
