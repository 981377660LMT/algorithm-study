// https://leetcode.cn/problems/memoize/

// The class has three public methods:

// set(key, value, duration): accepts an integer key, an integer value,
// and a duration in milliseconds. Once the duration has elapsed, the key should be inaccessible.
//  The method should return true if the same un-expired key already exists and false otherwise.
//   Both the value and duration should be overwritten if the key already exists.

// get(key): if an un-expired key exists, it should return the associated value.
//  Otherwise it should return -1.

// count(): returns the count of un-expired keys.

class TimeLimitedCache {
  private readonly _cache: Map<number, [value: number, expired: number]> = new Map()

  set(key: number, value: number, duration: number): boolean {
    this._expire()
    const res = this._cache.has(key)
    const now = Date.now()
    this._cache.set(key, [value, now + duration])
    return res
  }

  get(key: number): number {
    this._expire()
    if (!this._cache.has(key)) return -1
    return this._cache.get(key)![0]
  }

  count(): number {
    this._expire()
    return this._cache.size
  }

  /**
   * 朴素的实现，每次都遍历一遍，时间复杂度为 O(n).
   * 比较好的方法是维护一个优先队列，每次只需要取出队首元素即可(RemovableHeap)。
   */
  private _expire(): void {
    const now = Date.now()
    this._cache.forEach((value, key) => {
      if (value[1] <= now) this._cache.delete(key)
    })
  }
}

/**
 * Your TimeLimitedCache object will be instantiated and called as such:
 * var obj = new TimeLimitedCache()
 * obj.set(1, 42, 1000); // false
 * obj.get(1) // 42
 * obj.count() // 1
 */

export {}
