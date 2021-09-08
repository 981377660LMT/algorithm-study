const repeatFunc = repeat(console.log, 4, 1000, true)
repeatFunc('hellworld') //先立即打印一个hellworld，然后每个三秒打印三个hellworld

function repeat(this: any, func: Function, times: number, interval: number, immediate: boolean) {
  let count = 0
  const ctx = this

  // 有点像回溯
  function inner(...args: any[]) {
    count++
    if (count > times) return

    if (count === 1 && immediate) {
      func.call(ctx, ...args)
      inner.call(ctx, ...args)
      return
    }

    return setTimeout(() => {
      func.call(ctx, ...args)
      inner.call(ctx, ...args)
    }, interval)
  }

  return inner
}
