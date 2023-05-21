export {}

// 需要立即执行
// !应用场景：滚动事件，resize事件，鼠标移动事件（mousemove），DOM元素拖拽（mousemove），抢购疯狂点击（click）

type F = (...args: unknown[]) => void

function throttle(fn: F, t: number): F {
  let timer: ReturnType<typeof setTimeout> | null = null
  let preArgs: unknown[] | null = null

  return function dfs(...args) {
    if (timer) {
      preArgs = args
      return
    }

    fn(...args)
    timer = setTimeout(() => {
      timer = null
      if (preArgs) {
        dfs(...preArgs)
        preArgs = null
      }
    }, t)
  }
}
