interface ThrottleOption {
  leading: boolean // 是否立即执行
  trailing: boolean // 是否在冷却后执行
}

/**
 * @param {Function} func
 * @param {number} wait
 * 4. 手写throttle() 实际上是 {leading: true, trailing: true}的特殊情形。
 * leading:立即出发
 * trailing:如果之前被调用(有lastArgs) 则结束时多触发一次
 */
function throttle(
  func: Function,
  wait: number,
  option: ThrottleOption = { leading: true, trailing: true }
) {
  // your code here
  let timer: ReturnType<typeof setTimeout> | null = null
  let lastArgs: any = null
  const { leading, trailing } = option

  function setTimer(this: any) {
    if (lastArgs !== null && trailing) {
      func.apply(this, lastArgs)
      lastArgs = null
      timer = setTimeout(setTimer, wait)
    } else {
      timer = null
    }
  }

  return function (this: any, ...args: any[]) {
    if (!timer) {
      // 是否立即执行
      leading && func.apply(this, args)
      timer = setTimeout(setTimer, wait)
    } else {
      lastArgs = args
    }
  }
}

export {}
