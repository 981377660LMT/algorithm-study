export const raf = function (callback: () => void) {
  if (typeof requestAnimationFrame === typeof undefined) {
    const id = setTimeout(callback)
    return () => {
      clearTimeout(id)
    }
  }

  let rafId = 0

  rafId = requestAnimationFrame(callback)

  const clearRafTimeout = function () {
    cancelAnimationFrame(rafId)
  }

  return clearRafTimeout
}
