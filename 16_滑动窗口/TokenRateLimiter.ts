/**
 * TokenBucket 令牌桶流量控制算法实现.
 *
 * Token Bucket 是一种常用的流量控制算法，用于控制数据包发送的速率，确保网络不被过载。
 */
class TokenBucket {
  private readonly _refillRate: number // 令牌补充速率（每秒补充的令牌数）
  private readonly _capacity: number // 桶的容量
  private _numToken: number // 当前的令牌数量
  private _lastRefill: number // 上次补充令牌的时间戳

  /**
   * @param capacity 桶的最大容量
   * @param refillRate 每秒补充的令牌数
   */
  constructor(capacity: number, refillRate: number) {
    this._capacity = capacity
    this._refillRate = refillRate
    this._numToken = capacity
    this._lastRefill = Date.now()
  }

  /**
   * 尝试消费指定数量的令牌
   * @param count 需要消费的令牌数量
   * @returns 是否成功消费
   */
  allowed(count: number): boolean {
    this._refill()
    if (this._numToken >= count) {
      this._numToken -= count
      return true
    }
    return false
  }

  /**
   * 获取当前令牌数量
   * @returns 当前令牌数量
   */
  countToken(): number {
    this._refill()
    return this._numToken
  }

  /**
   * 补充令牌
   */
  private _refill() {
    const now = Date.now()
    const elapsed = (now - this._lastRefill) / 1000 // 以秒为单位
    const refillTokens = elapsed * this._refillRate
    if (refillTokens > 0) {
      this._numToken = Math.min(this._capacity, this._numToken + refillTokens)
      this._lastRefill = now
    }
  }
}

// 示例用法
const bucket = new TokenBucket(10, 5) // 容量10，每秒补充5个令牌

// 模拟发送数据
function sendData(dataSize: number) {
  if (bucket.allowed(dataSize)) {
    console.log(`发送数据: ${dataSize} 个令牌`)
  } else {
    console.log(`发送失败，当前令牌数: ${bucket.countToken().toFixed(2)}`)
  }
}

// 定时尝试发送数据
setInterval(() => {
  sendData(9)
}, 1000)
