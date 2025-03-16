/**
 * @param {Function} func
 * @param {number} wait
 */
function throttle(func: (...args: any[]) => void, wait: number) {
  let timer: ReturnType<typeof setTimeout> | null = null
  let lastArgs: any = null

  // 此处是取最后传入的参数执行
  return function (this: any, ...args: any[]) {
    if (timer) {
      lastArgs = args
    } else {
      // 立即执行
      func.apply(this, args)
      timer = setTimeout(() => {
        // 冷却后也执行
        lastArgs && func.apply(this, lastArgs)
        timer = null
      }, wait)
    }
  }
}
