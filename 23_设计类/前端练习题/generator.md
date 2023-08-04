`Generator<T, TReturn, TNext>` 是 TypeScript 中的一个泛型接口，用于描述生成器对象。它继承自 `Iterator<T, TReturn, TNext>` 接口。泛型参数具有以下含义：

- `T`：生成器产生的值的类型。
- `TReturn`：生成器函数中 `return` 语句返回的值的类型。
- `TNext`：在调用生成器的 `next()` 方法时传递的参数的类型。

`Generator` 接口定义了以下方法：

1. `next(...args: [] | [TNext]): IteratorResult<T, TReturn>`：通过传递可选参数 `TNext` 类型的值来恢复生成器的执行。返回一个 `IteratorResult<T, TReturn>` 对象，其中 `value` 属性包含生成器产生的值（类型为 `T`），`done` 属性表示生成器是否已完成。

2. `return(value: TReturn): IteratorResult<T, TReturn>`：结束生成器的执行，并返回一个带有 `value` 属性（类型为 `TReturn`）的 `IteratorResult<T, TReturn>` 对象。`done` 属性将为 `true`。

3. `throw(e: any): IteratorResult<T, TReturn>`：向生成器抛出一个异常。返回一个 `IteratorResult<T, TReturn>` 对象，其中 `value` 属性包含生成器产生的值（类型为 `T`），`done` 属性表示生成器是否已完成。

以下是一个使用 `Generator` 接口的示例：

```typescript
function* generatorFunction(): Generator<number, string, boolean> {
  const shouldYieldThree = yield 1
  if (shouldYieldThree) {
    yield 3
  } else {
    yield 2
  }
  return 'Done'
}

const generator = generatorFunction()

console.log(generator.next()) // { value: 1, done: false }
console.log(generator.next(true)) // { value: 3, done: false }
console.log(generator.return('Early return')) // { value: "Early return", done: true }
console.log(generator.throw(new Error('An error occurred'))) // 抛出异常
```

在这个示例中，`T` 是 `number`，`TReturn` 是 `string`，`TNext` 是 `boolean`。`next()` 方法接受一个 `boolean` 类型的参数，用于决定生成器函数中哪个 `yield` 语句将被执行。`return()` 方法接受一个 `string` 类型的参数，用于提前结束生成器函数的执行。`throw()` 方法接受一个异常对象，用于在生成器函数中抛出异常。
