/**
 * @param {Function} func
 * @param {number} wait
 */
function throttle(func: Function, wait: number) {
  // your code here
  let timer: NodeJS.Timer | null = null
  let lastArgs: any = null

  // 此处是取最后传入参数的节流类型
  return function (this: any, ...args: any[]) {
    if (timer) {
      lastArgs = args
    } else {
      // 是否立即执行
      func.apply(this, args)
      timer = setTimeout(() => {
        lastArgs && func.apply(this, lastArgs)
        timer = null
      }, wait)
    }
  }
}

if (require.main === module) {
}

// throttle在定时器里面清
// debounce在定时器外面清
