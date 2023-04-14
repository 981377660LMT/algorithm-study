/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-this-alias */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable eqeqeq */

import { Suspension } from './Suspension'

type Cell<C> = readonly [resolved: C, next: Stream<C>] | undefined | null

/**
 * 惰性求值的流。
 */
class Stream<S> extends Suspension<Cell<S>> {
  /**
   * 连接两个流。
   */
  static concat<T>(x: Stream<T>, y: Stream<T>): Stream<T> {
    return new Stream(() => {
      if (x.empty()) {
        return y.resolve()
      }
      return [x.top()!, Stream.concat(x.pop()!, y)]
    })
  }

  constructor()
  constructor(cell: Cell<S> | (() => Cell<S>))
  constructor(value: S, stream: Stream<S>)
  constructor(arg1?: S | Cell<S> | (() => Cell<S>), arg2?: Stream<S>) {
    const n = arguments.length // 实参个数
    if (n === 2) {
      super([arg1 as S, arg2 as Stream<S>])
    } else if (n === 1) {
      super(arg1 as Cell<S> | (() => Cell<S>))
    } else {
      super()
    }
  }

  empty(): boolean {
    return !this.resolve()
  }

  /**
   * 返回流的第一个元素。
   */
  top(): S | undefined {
    const res = this.resolve()
    return res ? res[0] : undefined
  }

  /**
   * 返回一个流，该流将第一个元素弹出。
   */
  pop(): Stream<S> | undefined {
    const res = this.resolve()
    return res ? res[1] : undefined
  }

  /**
   * 返回一个流，该流将元素`x`推送到前面。
   */
  push(x: S): Stream<S> {
    return new Stream(x, this)
  }

  /**
   * 返回流的反向流。
   */
  reverse(): Stream<S> {
    return new Stream(() => {
      let x: Stream<S> = this
      let res = new Stream<S>()
      while (!x.empty()) {
        res = res.push(x.top()!)
        x = x.pop()!
      }
      return res.resolve()
    })
  }

  override toString(): string {
    let x: Stream<S> = this
    const res: S[] = []
    while (!x.empty()) {
      res.push(x.top()!)
      x = x.pop()!
    }
    res.reverse()
    return `Stream{${res.join(', ')}}`
  }
}

export { Stream }

if (require.main === module) {
  const stream = new Stream<number>()
  const stream2 = stream.push(1).push(2).push(3)
  console.log(stream2.top())
  const reversed = stream2.reverse()
  console.log(reversed.top(), 111)
  console.log(reversed.toString())
}
