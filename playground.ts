export {}

/**
 * const throttled = throttle(console.log, 100);
 * throttled("log"); // logged immediately.
 * throttled("log"); // logged at t=100ms.
 */

type F = (...args: any[]) => void

function throttle(fn: F, t: number): F {
  let timer: ReturnType<typeof setTimeout> | null = null
  let store: any[] = []

  return function (...args) {
    if (timer) {
      store = args
      return
    }
    timer = setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, t)
  }
}

function jsonToMatrix(arr: any[]): (string | number | boolean | null)[] {}
