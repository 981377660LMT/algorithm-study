// 生成器函数

export {}

/**
 * T: yield 的类型.
 *
 * TReturn: 外面调用 gen.return(arg) 时，arg 的类型; 生成器return时的类型.
 *
 * TNext: 外面调用 gen.next(arg) 时，arg 的类型; arg 在生成器内部可以通过 `const arg = yield xxx` 获取.
 */
function* demo(): Generator<boolean, number, string> {
  while (true) {
    try {
      const a = yield true
      console.log(a)
      yield false
      return 111
      yield false
    } catch (error) {
      console.log()
    }
  }
}

const gen = demo()
gen.next()
// 向还未done的生成器抛出异常，并恢复生成器的执行.
// 如果生成器已经done了,会抛出异常.
gen.throw(new Error('error'))
console.log(gen.return(999))
console.log(gen.next())
console.log(gen.next())

export {}
