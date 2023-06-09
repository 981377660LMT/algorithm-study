function cancellable(fn: Function, args: any[], t: number): Function {
  fn(...args)
  const timer = setInterval(() => {
    fn(...args)
  }, t)
  return () => clearInterval(timer)
}

// !使用setTimeout实现
function cancellable2(fn: Function, args: any[], t: number): Function {
  let ok = true

  function run() {
    if (!ok) return
    fn(...args)
    setTimeout(() => {
      run()
    }, t)
  }

  run()

  return () => {
    ok = false
  }
}

/**
 *  const result = []
 *
 *  const fn = (x) => x * 2
 *  const args = [4], t = 20, cancelT = 110
 *
 *  const log = (...argsArr) => {
 *      result.push(fn(...argsArr))
 *  }
 *
 *  const cancel = cancellable(fn, args, t);
 *
 *  setTimeout(() => {
 *     cancel()
 *     console.log(result) // [
 *                         //      {"time":0,"returned":8},
 *                         //      {"time":20,"returned":8},
 *                         //      {"time":40,"returned":8},
 *                         //      {"time":60,"returned":8},
 *                         //      {"time":80,"returned":8},
 *                         //      {"time":100,"returned":8}
 *                         //  ]
 *  }, cancelT)
 */

export {}
