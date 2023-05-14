// 防抖
// 应用场景：搜索框输入，停止输入后，再发送请求

type F = (...p: unknown[]) => unknown

function debounce(fn: F, ms: number): F {
  let timer: ReturnType<typeof setTimeout> | null
  return function (this: unknown, ...args) {
    if (timer != null) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(this, args)
    }, ms)
  }
}

/**
 * const log = debounce(console.log, 100);
 * log('Hello'); // cancelled
 * log('Hello'); // cancelled
 * log('Hello'); // Logged at t=100ms
 */

export {}
