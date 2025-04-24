/**
 * 使用raf模拟 setInterval
 * @param callback
 * @param delay
 * @returns ()=>void
 */
export const setRafInterval = function (callback: () => void, delay = 0) {
  if (typeof requestAnimationFrame === typeof undefined) {
    const id = setTimeout(callback, delay)
    return () => {
      clearTimeout(id)
    }
  }

  let rafId = 0

  let startTime = new Date().getTime()

  const loop = () => {
    const current = new Date().getTime()
    if (current - startTime >= delay) {
      callback()
      startTime = new Date().getTime()
    }
    rafId = requestAnimationFrame(loop)
  }
  rafId = requestAnimationFrame(loop)

  const clearRafTimeout = function () {
    cancelAnimationFrame(rafId)
  }

  return clearRafTimeout
}
