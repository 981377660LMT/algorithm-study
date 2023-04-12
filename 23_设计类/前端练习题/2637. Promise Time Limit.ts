type Fn = (...params: any[]) => Promise<any>

function timeLimit(fn: Fn, ms: number): Fn {
  return async function (...args) {
    return Promise.race([
      fn(...args),
      new Promise((_, reject) => {
        // eslint-disable-next-line prefer-promise-reject-errors
        setTimeout(() => reject('Time Limit Exceeded'), ms)
      })
    ])
  }
}

/**
 * const limited = timeLimit((t) => new Promise(res => setTimeout(res, t)), 100);
 * limited(150).catch(console.log) // "Time Limit Exceeded" at t=100ms
 */

export {}
