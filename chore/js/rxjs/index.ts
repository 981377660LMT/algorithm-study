// Observable：表示一个可调用的未来值或事件的集合。
// Observer：一个回调函数的集合，知道如何监听由Observable提供的值。
// Subscription：表示Observable的执行，主要用于取消Observable的执行。
// Operators：纯函数，使得处理集合如数组或事件易于处理。
// Subject：等同于EventEmitter，并且是将值或事件多路推送给多个Observer的唯一方式

//
// 在函数式编程（Functional Programming, FP）中，"lift"是一个重要的概念，
// !它指的是将一个普通的函数转换（或提升，lifting）为一个能够作用于更高级别抽象的函数。
// 这种转换使得原本只能操作简单类型的函数能够操作复杂的数据结构，如容器（Container）、函子（Functor）、单子（Monad）等。
// eg:
const f = (x: number) => x + 1
const liftedF = (arr: number[]) => arr.map(f)
console.log(liftedF([1, 2, 3, 4, 5]))

export {}
