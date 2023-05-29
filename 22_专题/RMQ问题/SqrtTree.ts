/**
 * !理论O(nloglogn),但是没有st表快.
 * 优点在于维护的幺半群不需要满足幂等性(例如加法群不具有幂等性).
 */
class SqrtTree<E> {
  private readonly _nums: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _layerLg: number[]
  private readonly _onLayer: number[]
  private readonly _pref: E[][]
  private readonly _suf: E[][]
  private readonly _btwn: E[][]

  constructor(nums: E[], e: () => E, op: (a: E, b: E) => E) {
    this._nums = nums
    this._e = e
    this._op = op

    const n = nums.length
    let lg = 0
    while (1 << lg < n) lg++

    const onLayer: number[] = Array(lg + 1).fill(0)
    const layerLog: number[] = []
    let nLayer = 0
    for (let i = lg; i > 1; i = (i + 1) >> 1) {
      onLayer[i] = nLayer++
      layerLog.push(i)
    }

    for (let i = lg - 1; ~i; i--) onLayer[i] = Math.max(onLayer[i], onLayer[i + 1])
    const pref: E[][] = Array(nLayer)
    const suf: E[][] = Array(nLayer)
    const btwn: E[][] = Array(nLayer)
    for (let i = 0; i < nLayer; i++) {
      pref[i] = Array(n).fill(0)
      suf[i] = Array(n).fill(0)
      btwn[i] = Array(1 << lg).fill(0)
    }

    for (let layer = 0; layer < nLayer; layer++) {
      const prevBSz = 1 << layerLog[layer]
      const bSz = 1 << ((layerLog[layer] + 1) >>> 1)
      const bCnt = 1 << (layerLog[layer] >>> 1)
      for (let l = 0; l < n; l += prevBSz) {
        const r = Math.min(l + prevBSz, n)
        for (let a = l; a < r; a += bSz) {
          const b = Math.min(a + bSz, r)
          pref[layer][a] = nums[a]
          for (let i = a + 1; i < b; i++) pref[layer][i] = op(pref[layer][i - 1], nums[i])
          suf[layer][b - 1] = nums[b - 1]
          for (let i = b - 2; i >= a; i--) suf[layer][i] = op(nums[i], suf[layer][i + 1])
        }
        for (let i = 0; i < bCnt && l + i * bSz < n; i++) {
          let val = suf[layer][l + i * bSz]
          btwn[layer][l + i * bCnt + i] = val
          for (let j = i + 1; j < bCnt && l + j * bSz < n; j++) {
            val = op(val, suf[layer][l + j * bSz])
            btwn[layer][l + i * bCnt + j] = val
          }
        }
      }
    }

    this._layerLg = layerLog
    this._onLayer = onLayer
    this._pref = pref
    this._suf = suf
    this._btwn = btwn
  }

  /**
   * [l, r).
   * 0 <= l <= r <= n.
   */
  query(l: number, r: number): E {
    r--
    if (l > r) return this._e()
    if (l === r) return this._nums[l]
    if (l + 1 === r) return this._op(this._nums[l], this._nums[r])
    const layer = this._onLayer[32 - Math.clz32(l ^ r)]
    const bSz = 1 << ((this._layerLg[layer] + 1) >>> 1)
    const bCnt = 1 << (this._layerLg[layer] >>> 1)
    const a = (l >> this._layerLg[layer]) << this._layerLg[layer]
    const lBlock = (((l - a) / bSz) | 0) + 1
    const rBlock = (((r - a) / bSz) | 0) - 1
    let val = this._suf[layer][l]
    if (lBlock <= rBlock) val = this._op(val, this._btwn[layer][a + lBlock * bCnt + rBlock])
    val = this._op(val, this._pref[layer][r])
    return val
  }
}

export { SqrtTree }

if (require.main === module) {
  const st = new SqrtTree(
    [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    () => 0,
    (a, b) => a + b
  )
  console.log(st.query(0, 10), st.query(0, 0))
}
