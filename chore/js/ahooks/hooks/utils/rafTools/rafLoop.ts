/**
 * 创建一个ref循环
 */
export const createRafLooper = (callback?: () => void) => {
  let rafId = 0
  const loop = () => {
    callback?.()
    rafId = requestAnimationFrame(loop)
  }
  rafId = requestAnimationFrame(loop)

  return function () {
    cancelAnimationFrame(rafId)
  }
}
