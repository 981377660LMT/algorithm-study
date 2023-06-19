/* eslint-disable implicit-arrow-linebreak */

/**
 * 函子是一个能把值装(value)在里面，通过传入函数(map)来变换容器内容的`容器`.
 */
interface IFunctor<T> {
  readonly value: T
  map<U>(f: (x: T) => U): IFunctor<U>
}

/**
 * Pointed 函子的构造器.
 */
interface IPointedConstructor {
  new <V>(value: V): IFunctor<V>
  of<V>(value: V): IFunctor<V>
}

/**
 * 应用函子.
 * 实现了ap方法的Pointed函子.
 * Ap 函子可以让不同的函子可以相互应用，能够把一种函子的函数值应用到另一个函子的值上。
 */
interface IApplicative<T> extends IFunctor<T> {
  /**
   * ap 方法接收另一个函子作为参数，`调用其 map 方法`返回一个新的函子.
   */
  ap<U>(f: IFunctor<U>): IApplicative<U>
}

/**
 * 自函子.
 * 一种特殊的函子，其值和变形关系都是同一种类型.
 * !总是返回一个单层的函子，避免出现嵌套的情况.因此有一个`flatMap`方法.
 * Monad 函子的应用，就是实现 IO（输入、输出）操作。
 * @see {@link https://zhuanlan.zhihu.com/p/269513973}
 */
interface IEndoFunctor<T> extends IFunctor<T> {
  flatMap<U>(f: (x: T) => IEndoFunctor<U>): IEndoFunctor<U>
}

/**
 * 幺半群.
 */
interface IMonoid<E> {
  e: () => E
  op: (a: E, b: E) => E
}

/**
 * 单子.自函子上的幺半群.
 * @example
 * !Promise 是一个单子.
 * - 满足幺半群性质：
 *   - Promise.resolve(x=>x) => e()
 *   - Promise.then(x=>f(x).then(g)) => op(f,g)
 * - 满足自函子性质：
 *   - Promise.then(f) => flatMap (总是返回Promise本身,而不是嵌套的Promise)
 * @see {@link https://zhuanlan.zhihu.com/p/32734492}
 */
interface IMonad<T> extends IEndoFunctor<T>, IMonoid<T> {}

/**
 * IO 函子.是一种特殊的单子.
 * !内部的value是一个函数，这个函数是惰性的，只有调用了run方法，才会执行.
 * I/O 是不纯的操作，普通的函数式编程没法做，就需要把 IO 操作写成 Monad 函子来完成。
 * @see {@link https://zhuanlan.zhihu.com/p/56810671}
 */
interface IO<T> extends Omit<IMonad<T>, 'value'> {
  value: () => T
  run: () => T
}

//
//
//
// !应用:模拟java8的Stream
/**
 * 惰性求值的流.类似IO函子.
 * @alias IO
 */
class Stream<T> {
  /**
   * 将一个值包装为流.
   * @alias Pointed/Promise.resolve
   */
  static of<V>(value: V | Stream<V>): Stream<V> {
    if (value instanceof Stream) {
      return new Stream(() => value.run())
    }
    return new Stream(() => value)
  }

  private readonly _value: () => T

  constructor(effect: () => T) {
    this._value = effect
  }

  /**
   * @alias bind/Promise.prototype.then
   */
  map<U>(f: (x: T) => U): Stream<U> {
    return new Stream(() => f(this._value()))
  }

  /**
   * @alias bind/Promise.prototype.then
   */
  flatMap<U>(f: (x: T) => Stream<U>): Stream<U> {
    return new Stream(() => f(this._value()).run())
  }

  /**
   * 流的终止操作.
   * @alias join
   */
  run(): T {
    return this._value()
  }
}

export {}

if (require.main === module) {
  // !1.Stream函子的应用
  const readFile = (name: string) =>
    new Stream(() => {
      console.log("read file's content")
      return `请求并读取文件${name}的内容`
    })

  const format = (content: string) =>
    new Stream(() => {
      console.log('format content')
      return `格式化${content}`
    })

  const writeFile = (content: string) =>
    new Stream(() => {
      console.log('write file')
      return `写入${content}`
    })

  // 并没有执行
  const readAndWrite = readFile('test.txt').flatMap(format).flatMap(writeFile)
  readAndWrite.run()

  // !2.Promise就是一个单子.
  Promise.resolve(Promise.resolve(Promise.resolve(Promise.resolve(1)))).then(console.log)
  Stream.of(Stream.of(Stream.of(Stream.of(1))))
    .map(console.log)
    .run()
}
