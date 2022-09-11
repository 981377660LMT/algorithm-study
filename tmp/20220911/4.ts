// TODO 更新线段树模板
// 关注:不合理的区间的边界 tree的默认返回值

/**
 * @description 线段树区间最大值RMQ
 *
 * 如果查询超出范围 返回0
 */
class RMQSegmentTree {
  private readonly tree: number[]
  private readonly lazyValue: number[]
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   * @param size 区间右边界
   */
  constructor(size: number) {
    this.size = size
    this.tree = Array(size << 2).fill(0)
    this.lazyValue = Array(size << 2).fill(0)
    this.isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    if (l < 1) l = 1
    if (r > this.size) r = this.size
    if (l > r) return 0 // !超出范围返回0
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, target: number): void {
    if (l < 1) l = 1
    if (r > this.size) r = this.size
    if (l > r) return
    this._update(1, l, r, 1, this.size, target)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = 0 // !默认的最小值为0
    if (L <= mid) res = Math.max(res, this._query(rt << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((rt << 1) | 1, L, R, mid + 1, r))

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: number): void {
    if (L <= l && r <= R) {
      this.lazyValue[rt] = Math.max(this.lazyValue[rt], target)
      this.tree[rt] = Math.max(this.tree[rt], target)
      this.isLazy[rt] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this.tree[rt] = Math.max(this.tree[rt << 1], this.tree[(rt << 1) | 1])
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this.isLazy[rt]) {
      const target = this.lazyValue[rt]
      this.lazyValue[rt << 1] = Math.max(this.lazyValue[rt << 1], target)
      this.lazyValue[(rt << 1) | 1] = Math.max(this.lazyValue[(rt << 1) | 1], target)
      this.isLazy[rt << 1] = 1

      this.tree[rt << 1] = Math.max(this.tree[rt << 1], target)
      this.tree[(rt << 1) | 1] = Math.max(this.tree[(rt << 1) | 1], target)
      this.isLazy[(rt << 1) | 1] = 1

      this.lazyValue[rt] = 0
      this.isLazy[rt] = 0
    }
  }
}

// The subsequence is strictly increasing and
// The difference between adjacent elements in the subsequence is at most k.
// Return the length of the longest subsequence that meets the requirements.
// !注意到值域不对劲(只有1e5) 所以dp考虑以每个数为结尾
function lengthOfLIS(nums: number[], k: number): number {
  const max = Math.max(...nums)
  const dp = new RMQSegmentTree(max + 10)
  const n = nums.length
  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const maxLen = dp.query(num - k, num - 1) + 1 // !注意是严格递增
    dp.update(num, num, maxLen)
  }

  return dp.queryAll()
}

console.log(lengthOfLIS([4, 2, 1, 4, 3, 4, 5, 8, 15], 3))
console.log(lengthOfLIS([1, 100, 500, 100000, 100000], 100000))

export {}
