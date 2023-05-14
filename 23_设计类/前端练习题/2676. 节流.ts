export {}

// 需要立即执行
// !应用场景：滚动事件，resize事件，鼠标移动事件（mousemove），DOM元素拖拽（mousemove），抢购疯狂点击（click）

type F = (...args: any[]) => void

function throttle(fn: F, t: number): F {
  let timer: ReturnType<typeof setTimeout> | null
  let preArgs: any

  return function (...args) {
    if (timer) {
      preArgs = args
      return
    }
    fn.apply(this, args)
    timer = setTimeout(() => {
      if (preArgs) {
        fn.apply(this, preArgs)
        preArgs = null
      }
      timer = null
    }, t)
  }
}
