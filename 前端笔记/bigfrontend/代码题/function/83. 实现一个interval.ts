// 请实现一个mySetInterval(a, b)，这个和window.setInterval()不太一样，
// 时间间隔不是固定的，而是一个线性函数越来越大 period = a + b * count。

interface TimerRef {
  timer?: NodeJS.Timer | undefined
}
/**
 * @param {Function} func
 * @param {number} delay
 * @param {number} period
 * @return {number}
 */
function mySetInterval(func: Function, delay: number, period: number) {
  const timerRef: TimerRef = { timer: undefined }

  function inner(count: number) {
    timerRef.timer = setTimeout(() => {
      func()
      inner(count + 1)
    }, delay + period * count)
  }
  inner(0)

  return timerRef
}

/**
 * @param { number } id
 */
function myClearInterval(timerRef: TimerRef) {
  timerRef.timer && clearTimeout(timerRef.timer)
}

if (require.main === module) {
  let prev = Date.now()
  const func = () => {
    const now = Date.now()
    console.log('roughly ', Date.now() - prev)
    prev = now
  }

  const id = mySetInterval(func, 100, 200)
  // 100左右，100 + 200 * 0
  // 400左右，100 + 200 * 1
  // 900左右，100 + 200 * 2
  // 1600左右，100 + 200 * 3
  // ....

  // myClearInterval(id) // 停止interval
}

export {}
