// /* eslint-disable no-console */
// /* eslint-disable @typescript-eslint/no-this-alias */
// /* eslint-disable @typescript-eslint/no-non-null-assertion */
// /* eslint-disable eqeqeq */

// class Stream<S> {
//   static concat<T>(x: Stream<T>, y: Stream<T>): Stream<T> {
//     if (!x._next) {
//       return new Stream(y._value, y._next)
//     }
//     return new Stream(x._value!, Stream.concat(x._next!, y))
//   }

//   private readonly _value: S | undefined
//   private readonly _next: Stream<S> | undefined

//   constructor(value?: S, next?: Stream<S>) {
//     this._value = value
//     this._next = next
//   }

//   empty(): boolean {
//     return !this._next
//   }

//   top(): S | undefined {
//     return this._value
//   }

//   pop(): Stream<S> | undefined {
//     return this._next
//   }

//   push(x: S): Stream<S> {
//     return new Stream(x, this)
//   }

//   reverse(): Stream<S> {
//     let x: Stream<S> = this
//     let res = new Stream<S>()
//     while (x._next) {
//       res = res.push(x._value!)
//       x = x._next!
//     }
//     return new Stream(res._value, res._next)
//   }

//   toString(): string {
//     let x: Stream<S> = this
//     const res: S[] = []
//     while (x._next) {
//       res.push(x._value!)
//       x = x._next!
//     }
//     res.reverse()
//     return `Link{${res.join(', ')}}`
//   }
// }

// export {}

// if (require.main === module) {
//   const link = new Stream<number>()
//   const link2 = link.push(1).push(2).push(3)
//   console.log(link2.top())
//   const reversed = link2.reverse()
//   console.log(reversed.top(), 111)
//   console.log(reversed.toString())
// }
