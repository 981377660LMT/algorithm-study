/**
 * 使用raf模拟 setTimeout
 * @param callback
 * @param delay
 * @returns ()=>void
 */
export const setRafTimeout = function (callback: () => void, delay = 0) {
  if (typeof requestAnimationFrame === typeof undefined) {
    const id = setTimeout(callback, delay)
    return () => {
      clearTimeout(id)
    }
  }

  let rafId = 0

  const startTime = new Date().getTime()

  const loop = () => {
    const current = new Date().getTime()
    if (current - startTime >= delay) {
      callback()
    } else {
      rafId = requestAnimationFrame(loop)
    }
  }
  rafId = requestAnimationFrame(loop)

  const clearRafTimeout = function () {
    cancelAnimationFrame(rafId)
  }

  return clearRafTimeout
}
