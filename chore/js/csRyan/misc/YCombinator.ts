// 假设某个编程语言不支持函数调用本身，如何实现递归
// Y组合子
// !思路是，虽然函数不能直接引用自己，但是可以使用参数！
// 设有一个函数 f，想要用递归来定义：
// f(x) = ... // 需要引用 f
// 通过 Y 组合子，可以写成：
// !Y(f)(x) = f(Y(f)(x))  // 这允许 f 在其内部引用自身。

{
  function factorial(n: number): number {
    return _factorial(_factorial, n)
  }

  function _factorial(self: (self: typeof _factorial, n: number) => number, n: number): number {
    if (n === 0) return 1
    return n * self(self, n - 1)
  }

  console.log(factorial(5)) // 120
}

// YCombinator

type Calculator<T, R> = (arg: T) => R
function YCombinator<T, R>(f: (getSelf: Calculator<T, R>) => Calculator<T, R>): Calculator<T, R> {
  return (arg: T): R => f(YCombinator(f))(arg)
}

const _fac =
  (getSelf: (arg: number) => number) =>
  (n: number): number => {
    if (n === 0) return 1
    return n * getSelf(n - 1)
  }

const factorial = YCombinator(_fac)

console.log(factorial(5)) // 120

export {}
