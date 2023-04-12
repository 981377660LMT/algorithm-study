// Not all C++ constructs are directly translatable to TypeScript,
// but here is a TypeScript approximation of the original code:

type Function<T> = () => T

class Suspension<T> {
  private value: T | undefined
  private lazy: Function<T> | undefined

  constructor(x: T)
  constructor(f: Function<T>)
  constructor(arg: T | Function<T>) {
    if (typeof arg === 'function') {
      this.lazy = arg
    } else {
      this.value = arg
    }
  }

  force(): T {
    if (this.lazy !== undefined) {
      this.value = this.lazy()
      this.lazy = undefined
    }
    if (this.value === undefined) {
      throw new Error('Suspension was not properly initialized')
    }
    return this.value
  }
}

type Pair<T> = [T, Stream<T>]

class Stream<T> extends Suspension<Pair<T> | undefined> {
  constructor(x?: T, s?: Stream<T>) {
    super([[x, s]])
  }

  empty(): boolean {
    return this.force() === undefined
  }

  top(): T {
    const [x] = this.force()!
    if (x === undefined) {
      throw new Error('Stream is empty')
    }
    return x
  }

  push(x: T): Stream<T> {
    return new Stream(x, this)
  }

  pop(): Stream<T> {
    const [x, s] = this.force()!
    if (x === undefined) {
      throw new Error('Stream is empty')
    }
    return s
  }

  reverse(): Stream<T> {
    const f = () => {
      let ret = new Stream<T>()
      let x = this
      while (!x.empty()) {
        ret = ret.push(x.top())
        x = x.pop()
      }
      return ret.force()
    }
    return new Stream(f)
  }

  static concat<T>(l: Stream<T>, r: Stream<T>): Stream<T> {
    const f = () => {
      const [x, s] = l.force()!
      if (x === undefined) {
        return r.force()
      } else {
        return [[x[0], Stream.concat(x[1], r)], s].flat() as Pair<T> | undefined
      }
    }
    return new Stream(f)
  }
}
