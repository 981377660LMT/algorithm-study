import { Func } from '../../typings'

/**
 * @param {Function} func
 * @param {number} wait
 * @param {boolean} option.leading
 * @param {boolean} option.trailing
 * 6. 手写debounce() 实际上是 {leading: false, trailing: true}的特殊情形。
 * leading:开始时立即出发
 * trailing:如果之前没有被调用(!isCalled) 则结束时触发一次
 */
function debounce(func: Function, wait: number, option = { leading: false, trailing: true }): Func {
  let timer: NodeJS.Timer | null = null
  const { leading, trailing } = option

  return function (this: unknown, ...args) {
    let isCalled = false

    // 开始时触发吗
    if (!timer && leading) {
      isCalled = true
      func.apply(this, args)
    }

    timer && clearTimeout(timer)

    timer = setTimeout(() => {
      !isCalled && trailing && func.apply(this, args)
      timer = null
    }, wait)
  }
}
