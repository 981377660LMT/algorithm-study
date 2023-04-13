// TODO 没弄明白

function bootStrap<Args, Return>(
  wrapped: (...args: Args[]) => Generator<Return>,
  stack: Generator[] = []
): typeof wrapped {
  // eslint-disable-next-line func-names
  const gen = function (...args: Args[]): Generator<Return> {
    if (stack.length) {
      return wrapped(...args)
    }

    let to: Generator<Return> = wrapped(...args)
    while (true) {
      if (isGenerator(to)) {
        stack.push(to)
        to = to.next().value
      } else {
        stack.pop()
        if (!stack.length) {
          break
        }
        to = stack[stack.length - 1].next(to).value
      }
    }

    return to
  }

  return gen
}

function isGenerator(obj: any): obj is Generator {
  return obj && obj.constructor && obj.constructor.name === 'GeneratorFunction'
}

export {}

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  let fib = function* fibonacci(n: number, a = 0, b = 1): Generator<number> {
    if (n > 0) {
      yield a
      yield* fibonacci(n - 1, b, a + b)
    }
  }
  fib = bootStrap(fib)
  console.log('fib(10):', [...fib(10000)]) // 没什么用
}
