// 四个接口的类型，Consumer(消费型)、Supplier(供给型)、Predicate(判断型)与Function(转换型)
// 对应的抽象方法Consumer(accpet)、Supplier(get)、Predicate(test)与Function(apply)

/**
 * Consumer接口是一个消费型的接口，只要实现它的accept方法，就能作为消费者来输出信息。
 * lambda、方法引用都可以是一个Consumer类型，因此他们可以作为forEach的参数，用来协助Stream输出信息。
 */
interface IConsumer<T> {
  accpet(t: T): void
  andThen<R extends T>(after: IConsumer<R>): IConsumer<T>
}

/**
 * Supplier接口是一个供给型的接口，只要实现它的get方法，就能作为供给者来提供信息。
 */
interface ISupplier<T> {
  get(): T
}

/**
 * Predicate接口是一个判断型的接口，只要实现它的test方法，就能作为判断者来判断信息。
 */
interface IPredicate<T> {
  test(t: T): boolean
  and<R extends T>(other: IPredicate<R>): IPredicate<T>
  or<R extends T>(other: IPredicate<R>): IPredicate<T>
  negate(): IPredicate<T>
}

/**
 * Function接口是一个转换型的接口，只要实现它的apply方法，就能作为转换者来转换信息。
 */
interface IFunction<T, R> {
  apply(t: T): R
  compose<V>(before: IFunction<V, T>): IFunction<V, R>
  andThen<V>(after: IFunction<R, V>): IFunction<T, V>
  identity(): IFunction<T, T>
}

export {}
