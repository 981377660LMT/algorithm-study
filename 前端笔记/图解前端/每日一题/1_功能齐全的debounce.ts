function superDebounce(fn: (...args: any[]) => any, delay: number) {
  let timer: NodeJS.Timer | null = null

  // this的指向要跟原来函数一样
  function run(this: any, ...args: any[]) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.call(this, ...args)
    }, delay)
  }

  function cancel() {
    timer && clearTimeout(timer)
    timer = null
  }

  return { run, cancel }
}

const testFunc = () => {
  console.log('test')
  return 666
}
const debouncedTestFunc = superDebounce(testFunc, 500)
console.log(debouncedTestFunc.run())
// debouncedTestFunc.run()
