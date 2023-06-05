// 2715. 延迟执行的可取消函数
function cancellable(fn: Function, args: any[], t: number): Function {
  const timer = setTimeout(() => {
    fn(...args)
  }, t)
  return () => clearTimeout(timer)
}

/**
 *  const result = []
 *
 *  const fn = (x) => x * 5
 *  const args = [2], t = 20, cancelT = 50
 *
 *  const log = (...argsArr) => {
 *      result.push(fn(...argsArr))
 *  }
 *
 *  const cancel = cancellable(fn, args, t);
 *
 *  setTimeout(() => {
 *     cancel()
 *     console.log(result) // [{"time":20,"returned":10}]
 *  }, cancelT)
 */

export {}
