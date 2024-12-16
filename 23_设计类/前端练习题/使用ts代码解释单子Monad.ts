/* eslint-disable no-inner-declarations */
/* eslint-disable no-lone-blocks */

// 强类型版 Continuation Passing Style
// 不要把 Monad 想象的太玄乎、不要掉进范畴论，你的思考可以用在更务实的地方上。
//
// - Monad 的组成：
//   - 类型构造器：将普通类型转换为 Monad 类型。
//   - `of` 方法：将普通值包装到 Monad 中。
//   - `flatMap` 方法：将 Monad 中的值与一个返回 Monad 的函数组合。

// - Monad 的三大法则：
//   - 左单位元律、右单位元律、结合律，确保 Monad 的行为一致。

// - Maybe Monad 示例：
//   - 处理可能为空的值，通过链式调用自动传播错误（`Nothing`）。

// - 优势：
//   - 简化错误处理、提高代码可组合性、管理副作用。

abstract class Monad<T> {
  /** 将一个值包装到 Monad 中 */
  static of<U>(_value: U): Monad<U> {
    throw new Error('Not implemented')
  }

  /** 将 Monad 中的值与一个返回 Monad 的函数组合 */
  abstract flatMap<U>(fn: (value: T) => Monad<U>): Monad<U>

  /** 额外的 map 方法，便于链式调用 */
  map<U>(fn: (value: T) => U): Monad<U> {
    return this.flatMap(value => Monad.of(fn(value)))
  }
}

type Maybe<T> = Just<T> | Nothing

class Just<T> extends Monad<T> {
  static override of<U>(_value: U): Monad<U> {
    return new Just(_value)
  }

  readonly value: T

  constructor(value: T) {
    super()
    this.value = value
  }

  override flatMap<U>(fn: (value: T) => Monad<U>): Monad<U> {
    return fn(this.value)
  }
}

class Nothing extends Monad<never> {
  static override of<U>(): Monad<U> {
    return new Nothing()
  }

  // eslint-disable-next-line class-methods-use-this
  override flatMap<U>(): Monad<U> {
    return new Nothing()
  }
}

function of<T>(value: T): Maybe<T> {
  if (value === null || value === undefined) {
    return new Nothing()
  }
  return new Just(value)
}

{
  function safeParseNumber(str: string): Maybe<number> {
    const num = parseFloat(str)
    return Number.isNaN(num) ? new Nothing() : new Just(num)
  }

  const res = of('1')
    .flatMap(str => safeParseNumber(str))
    .flatMap(num => new Just(num + 1))

  if (res instanceof Just) {
    console.log(res.value)
  } else {
    console.log('Nothing')
  }
}

{
  // `Promise` 本身就是一个 Monad，因为它满足 Monad 的三大法则。我们可以利用 TypeScript 中的 `Promise` 来展示 Monad 的行为。
  // 示例函数：获取用户信息
  function getUser(userId: number): Promise<Maybe<{ id: number; name: string }>> {
    return new Promise(resolve => {
      // 模拟异步操作
      setTimeout(() => {
        if (userId === 1) {
          resolve(of({ id: 1, name: 'Alice' }))
        } else {
          resolve(new Nothing())
        }
      }, 1000)
    })
  }

  // 示例函数：获取用户的朋友
  function getFriends(user: { id: number; name: string }): Promise<Maybe<string[]>> {
    return new Promise(resolve => {
      setTimeout(() => {
        if (user.id === 1) {
          resolve(of(['Bob', 'Charlie']))
        } else {
          resolve(new Nothing())
        }
      }, 1000)
    })
  }

  // 链式调用
  getUser(1)
    .then(userMaybe => userMaybe.flatMap(user => getFriends(user)))
    .then(friendsMaybe => {
      if (friendsMaybe instanceof Just) {
        console.log('Friends:', friendsMaybe.value) // 输出: Friends: [ 'Bob', 'Charlie' ]
      } else {
        console.log('No friends found.')
      }
    })

  getUser(2)
    .then(userMaybe => userMaybe.flatMap(user => getFriends(user)))
    .then(friendsMaybe => {
      if (friendsMaybe instanceof Just) {
        console.log('Friends:', friendsMaybe.value)
      } else {
        console.log('No friends found.') // 输出: No friends found.
      }
    })
}

export {}
