// NoInfer 内置类型
//
// 在调用泛型函数时，TypeScript 能够从你传入的任何内容中推断出类型参数。
// !然而，一个挑战是，并不总是很清楚要推断的“最佳”类型是什么。这可能会导致 TypeScript 拒绝有效的调用，接受有问题的调用

declare function doSomething<T>(arg: T): void

// We can explicitly say that 'T' should be 'string'.
doSomething<string>('hello!')
// We can also just let the type of 'T' get inferred.
doSomething('hello!')

declare function createStreetLight<C extends string>(colors: C[], defaultColor?: C): void
createStreetLight(['red', 'yellow', 'green'], 'red')
createStreetLight(['red', 'yellow', 'green'], 'blue') // 不报错，不符合预期

// !人们目前处理这个问题的一种方法是添加一个单独的类型参数，该参数受现有类型参数的约束。
declare function createStreetLight2<C extends string, D extends C>(colors: C[], defaultColor?: D): void
createStreetLight2(['red', 'yellow', 'green'], 'blue') // 报错，符合预期

// 但是这有点麻烦。
// 使用 NoInfer，我们可以将 createStreetLight 重写为如下:
// 将 defaultColor 的类型排除在推理中意味着 “blue” 永远不会成为推理候选者，类型检查器可以拒绝它。
declare function createStreetLight3<C extends string>(colors: C[], defaultColor?: NoInfer<C>): void
createStreetLight3(['red', 'yellow', 'green'], 'blue')
