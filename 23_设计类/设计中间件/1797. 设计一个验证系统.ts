// 所有函数调用中，currentTime 的值 严格递增
type TokenId = string
type ExpiredTime = number

class AuthenticationManager {
  private timeToLive: number
  private lru: Map<TokenId, ExpiredTime>

  constructor(timeToLive: number) {
    this.timeToLive = timeToLive
    this.lru = new Map()
  }

  // 给定 tokenId ，在当前时间 currentTime 生成一个新的验证码
  generate(tokenId: string, currentTime: number): void {
    this.lru.set(tokenId, currentTime + this.timeToLive)
  }

  // 将给定 tokenId 且 未过期 的验证码在 currentTime 时刻更新
  // 如果给定 tokenId 对应的验证码不存在或已过期，请你忽略该操作，不会有任何更新操作发生。
  renew(tokenId: string, currentTime: number): void {
    if (!this.lru.has(tokenId)) return
    if (this.lru.get(tokenId)! > currentTime) {
      // 存活
      this.lru.delete(tokenId)
      this.lru.set(tokenId, currentTime + this.timeToLive)
    } else {
      // 过期
      this.lru.delete(tokenId)
    }
  }

  // 请返回在给定 currentTime 时刻，未过期 的验证码数目
  countUnexpiredTokens(currentTime: number): number {
    for (const [tokenId, expiredTime] of this.lru.entries()) {
      if (expiredTime > currentTime) break
      this.lru.delete(tokenId)
    }
    return this.lru.size
  }
}

export {}

// LinkedHashMap解法
// 使用LRU缓存的思想即可 可以改善countUnexpiredTokens的复杂度 (即提前break)
