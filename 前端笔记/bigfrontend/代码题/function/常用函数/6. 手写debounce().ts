/**
 * @param {Function} func
 * @param {number} wait
 */
function debounce(func: Function, wait: number) {
  let timer: ReturnType<typeof setTimeout> | null = null

  return function (this: any, ...args: any[]) {
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      func.apply(this, args)
    }, wait)
  }
}
