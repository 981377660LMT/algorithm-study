/** 简单限流. */
class SlidingWindowRateLimiter {
  private readonly _maxRequestCount: number
  private readonly _windowLengthMs: number
  private readonly _window: { requestCount: number; timeStamp: number }[] = []
  private _countInWindow = 0

  constructor(options: { maxRequestCount: number; windowLengthMs: number }) {
    this._maxRequestCount = options.maxRequestCount
    this._windowLengthMs = options.windowLengthMs
  }

  allowed(requestCount: number): boolean {
    const now = Date.now()
    let head = 0
    while (head < this._window.length && now - this._window[head].timeStamp > this._windowLengthMs) {
      this._countInWindow -= this._window[head].requestCount
      head++
    }
    this._window.splice(0, head)

    if (this._countInWindow + requestCount > this._maxRequestCount) {
      return false
    } else {
      this._window.push({ requestCount, timeStamp: now })
      this._countInWindow += requestCount
      return true
    }
  }

  clear(): void {
    this._window.length = 0
    this._countInWindow = 0
  }
}

export { SlidingWindowRateLimiter }

if (require.main === module) {
  const rateLimiter = new SlidingWindowRateLimiter({ maxRequestCount: 5, windowLengthMs: 2000 })
  console.log(rateLimiter.allowed(1))
  console.log(rateLimiter.allowed(1))
  console.log(rateLimiter.allowed(1))
  console.log(rateLimiter.allowed(1))
  console.log(rateLimiter.allowed(1))
  console.log(rateLimiter.allowed(1))

  setTimeout(() => {
    console.log(rateLimiter.allowed(1))
    console.log(rateLimiter.allowed(4))
    console.log(rateLimiter.allowed(1))
  }, 2000)
}
