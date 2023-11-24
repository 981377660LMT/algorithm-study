/* eslint-disable no-empty */
/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

import { useBlock } from '../../22_专题/离线查询/根号分治/SqrtDecomposition/useBlock'

class WaveletMatrixLikeOfflineDynamic {
  private readonly _nums: Uint32Array
  private readonly _valueToId: Map<number, number>
  private readonly _idToValue: number[]
  private readonly _sortedLength: number

  private readonly _belong1: Uint16Array
  private readonly _blockStart1: Uint32Array
  private readonly _blockEnd1: Uint32Array
  private readonly _blockCount1: number
  private readonly _belong2: Uint16Array
  private readonly _blockStart2: Uint32Array
  private readonly _blockEnd2: Uint32Array
  private readonly _blockCount2: number

  /** preSum1[i][j] 表示前i个块中有多少个数属于第j个值域块. */
  private readonly _preSum1: Uint32Array
  /** preSum2[i][j] 表示前i个块中有多少个数等于j. */
  private readonly _preSum2: Uint32Array
  /** fragmentCounter1[j] 临时保存散块内出现在第j个值域块中的数的个数. */
  private readonly _fragmentCounter1: Uint16Array
  /** fragmentCounter2[j] 临时保存散块内等于j的数的个数. */
  private readonly _fragmentCounter2: Uint16Array

  /**
   * 序列分块+值域分块模拟`wavelet matrix`.支持单点修改.
   * 单次操作时间复杂度`O(sqrt(n))`.
   *
   * @param initNums 初始化的序列.
   * @param allowedNums 允许出现的数, 包含修改和查询的数.
   * @param blockSize 序列分块的大小.为了减少空间复杂度`O(n*块数)`，可以把块的大小增大一些.默认为`4*sqrt(n)`.
   */
  constructor(initNums: number[], allowedNums: Iterable<number>, blockSize = 4 * (Math.sqrt(initNums.length + 1) | 0)) {
    const sorted = [...new Set(allowedNums)].sort((a, b) => a - b)
    const valueToId = new Map<number, number>()
    const idToValue = Array<number>(sorted.length)
    for (let i = 0; i < sorted.length; i++) {
      valueToId.set(sorted[i], i)
      idToValue[i] = sorted[i]
    }
    const newNums = new Uint32Array(initNums.length)
    for (let i = 0; i < initNums.length; i++) {
      newNums[i] = valueToId.get(initNums[i])!
    }

    const { belong: belong1, blockStart: blockStart1, blockEnd: blockEnd1, blockCount: blockCount1 } = useBlock(newNums.length, blockSize)
    const { belong: belong2, blockStart: blockStart2, blockEnd: blockEnd2, blockCount: blockCount2 } = useBlock(sorted.length)

    const preSum1 = new Uint32Array((blockCount1 + 1) * blockCount2)
    const preSum2 = new Uint32Array((blockCount1 + 1) * sorted.length)
    for (let bid = 0; bid < blockCount1; bid++) {
      preSum1.copyWithin((bid + 1) * blockCount2, bid * blockCount2, (bid + 1) * blockCount2)
      preSum2.copyWithin((bid + 1) * sorted.length, bid * sorted.length, (bid + 1) * sorted.length)
      for (let i = blockStart1[bid]; i < blockEnd1[bid]; i++) {
        const num = newNums[i]
        const vid = belong2[num]
        preSum1[(bid + 1) * blockCount2 + vid]++
        preSum2[(bid + 1) * sorted.length + num]++
      }
    }
    const counter1 = new Uint16Array(blockCount2)
    const counter2 = new Uint16Array(sorted.length)

    this._nums = newNums
    this._valueToId = valueToId
    this._idToValue = idToValue
    this._sortedLength = sorted.length

    this._belong1 = belong1
    this._blockStart1 = blockStart1
    this._blockEnd1 = blockEnd1
    this._blockCount1 = blockCount1
    this._belong2 = belong2
    this._blockStart2 = blockStart2
    this._blockEnd2 = blockEnd2
    this._blockCount2 = blockCount2

    this._preSum1 = preSum1
    this._preSum2 = preSum2
    this._fragmentCounter1 = counter1
    this._fragmentCounter2 = counter2
  }

  get(index: number): number | undefined {
    if (index < 0 || index >= this._nums.length) return undefined
    return this._idToValue[this._nums[index]]
  }

  set(index: number, target: number): void {
    if (index < 0 || index >= this._nums.length) return
    target = this._getId(target)
    if (this._nums[index] === target) return
    const preValue = this._nums[index]
    const curValue = target
    const preVid = this._belong2[preValue]
    const curVid = this._belong2[curValue]
    for (let bid = this._belong1[index]; bid < this._blockCount1; bid++) {
      this._preSum1[(bid + 1) * this._blockCount2 + preVid]--
      this._preSum1[(bid + 1) * this._blockCount2 + curVid]++
      this._preSum2[(bid + 1) * this._sortedLength + preValue]--
      this._preSum2[(bid + 1) * this._sortedLength + curValue]++
    }
    this._nums[index] = target
  }

  /**
   * 查询区间`[start, end)`中第`k(k>=0)`小的数.
   */
  kth(start: number, end: number, k: number): number | undefined {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end || end - start <= k) return undefined
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    if (bid1 === bid2) {
      this._updateFragment(start, end, 1)
      let todo = k + 1
      for (let vid = 0; vid < this._blockCount2; vid++) {
        if (todo > this._fragmentCounter1[vid]) {
          todo -= this._fragmentCounter1[vid]
          continue
        }
        for (let j = this._blockStart2[vid]; j < this._blockEnd2[vid]; j++) {
          todo -= this._fragmentCounter2[j]
          if (todo <= 0) {
            this._updateFragment(start, end, -1)
            return this._idToValue[j]
          }
        }
      }
    } else {
      this._updateFragment(start, this._blockEnd1[bid1], 1)
      this._updateFragment(this._blockStart1[bid2], end, 1)
      let todo = k + 1
      for (let vid = 0; vid < this._blockCount2; vid++) {
        const curCount =
          this._fragmentCounter1[vid] + this._preSum1[bid2 * this._blockCount2 + vid] - this._preSum1[(bid1 + 1) * this._blockCount2 + vid]
        if (todo > curCount) {
          todo -= curCount
          continue
        }
        for (let j = this._blockStart2[vid]; j < this._blockEnd2[vid]; j++) {
          const curCount =
            this._fragmentCounter2[j] + this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
          todo -= curCount
          if (todo <= 0) {
            this._updateFragment(start, this._blockEnd1[bid1], -1)
            this._updateFragment(this._blockStart1[bid2], end, -1)
            return this._idToValue[j]
          }
        }
      }
    }

    return undefined
  }

  floor(start: number, end: number, target: number): number | undefined {
    return this._prev(start, end, target, true)
  }

  ceiling(start: number, end: number, target: number): number | undefined {
    return this._next(start, end, target, true)
  }

  lower(start: number, end: number, target: number): number | undefined {
    return this._prev(start, end, target, false)
  }

  higher(start: number, end: number, target: number): number | undefined {
    return this._next(start, end, target, false)
  }

  count(start: number, end: number, target: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    target = this._getId(target)
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] === target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd1[bid1]; i++) res += +(this._nums[i] === target)
    for (let i = this._blockStart1[bid2]; i < end; i++) res += +(this._nums[i] === target)
    res += this._preSum2[bid2 * this._sortedLength + target] - this._preSum2[(bid1 + 1) * this._sortedLength + target]
    return res
  }

  /**
   * 查询区间`[start, end)`中在`[lower, upper)`范围内的数的个数.
   */
  countRange(start: number, end: number, lower: number, upper: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    if (lower >= upper) return 0
    lower = this._getId(lower)
    upper = this._getId(upper)

    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] >= lower && this._nums[i] < upper)
      return res
    }

    let res = 0
    for (let i = start; i < this._blockEnd1[bid1]; i++) res += +(this._nums[i] >= lower && this._nums[i] < upper)
    for (let i = this._blockStart1[bid2]; i < end; i++) res += +(this._nums[i] >= lower && this._nums[i] < upper)
    const vid1 = this._belong2[lower]
    const vid2 = this._belong2[upper]
    if (vid1 === vid2) {
      for (let j = lower; j < upper; j++) {
        res += this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
      }
    } else {
      for (let j = lower; j < this._blockEnd2[vid1]; j++) {
        res += this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
      }
      for (let j = vid1 + 1; j < vid2; j++) {
        res += this._preSum1[bid2 * this._blockCount2 + j] - this._preSum1[(bid1 + 1) * this._blockCount2 + j]
      }
      for (let j = this._blockStart2[vid2]; j < upper; j++) {
        res += this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
      }
    }
    return res
  }

  countFloor(start: number, end: number, target: number): number {
    return this._countPrev(start, end, target, true)
  }

  countCeiling(start: number, end: number, target: number): number {
    return this._countNext(start, end, target, true)
  }

  countLower(start: number, end: number, target: number): number {
    return this._countPrev(start, end, target, false)
  }

  countHigher(start: number, end: number, target: number): number {
    return this._countNext(start, end, target, false)
  }

  private _updateFragment(start: number, end: number, delta: number): void {
    for (let i = start; i < end; i++) {
      const num = this._nums[i]
      const vid = this._belong2[num]
      this._fragmentCounter1[vid] += delta
      this._fragmentCounter2[num] += delta
    }
  }

  private _getId(rawValue: number): number {
    const id = this._valueToId.get(rawValue)
    if (id == undefined) throw new Error(`rawValue ${rawValue} not in allNums`)
    return id
  }

  private _prev(start: number, end: number, target: number, inclusive: boolean): number | undefined {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return undefined
    target = this._getId(target)
    if (!inclusive && target === 0) return undefined
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    const vid = this._belong2[target]
    if (bid1 === bid2) {
      this._updateFragment(start, end, 1)
      let cand = inclusive ? target : target - 1
      const vStart = this._blockStart2[vid]
      for (; cand >= vStart && !this._fragmentCounter2[cand]; cand--) {}
      if (cand >= vStart) {
        this._updateFragment(start, end, -1)
        return this._idToValue[cand]
      }

      let candVid = vid - 1
      for (; candVid >= 0 && !this._fragmentCounter1[candVid]; candVid--) {}
      if (candVid === -1) {
        this._updateFragment(start, end, -1)
        return undefined
      }
      for (let j = this._blockEnd2[candVid] - 1; j >= this._blockStart2[candVid]; j--) {
        if (this._fragmentCounter2[j]) {
          this._updateFragment(start, end, -1)
          return this._idToValue[j]
        }
      }
    } else {
      this._updateFragment(start, this._blockEnd1[bid1], 1)
      this._updateFragment(this._blockStart1[bid2], end, 1)
      let cand = inclusive ? target : target - 1
      const vStart = this._blockStart2[vid]
      for (
        ;
        cand >= vStart &&
        !(this._fragmentCounter2[cand] + this._preSum2[bid2 * this._sortedLength + cand] - this._preSum2[(bid1 + 1) * this._sortedLength + cand]);
        cand--
      ) {}
      if (cand >= vStart) {
        this._updateFragment(start, this._blockEnd1[bid1], -1)
        this._updateFragment(this._blockStart1[bid2], end, -1)
        return this._idToValue[cand]
      }

      let candVid = vid - 1
      for (
        ;
        candVid >= 0 &&
        !(
          this._fragmentCounter1[candVid] +
          this._preSum1[bid2 * this._blockCount2 + candVid] -
          this._preSum1[(bid1 + 1) * this._blockCount2 + candVid]
        );
        candVid--
      ) {}
      if (candVid === -1) {
        this._updateFragment(start, this._blockEnd1[bid1], -1)
        this._updateFragment(this._blockStart1[bid2], end, -1)
        return undefined
      }
      for (let j = this._blockEnd2[candVid] - 1; j >= this._blockStart2[candVid]; j--) {
        if (this._fragmentCounter2[j] + this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]) {
          this._updateFragment(start, this._blockEnd1[bid1], -1)
          this._updateFragment(this._blockStart1[bid2], end, -1)
          return this._idToValue[j]
        }
      }
    }

    throw new Error('unreachable')
  }

  private _next(start: number, end: number, target: number, inclusive: boolean): number | undefined {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return undefined
    target = this._getId(target)
    if (!inclusive && target === this._sortedLength - 1) return undefined
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    const vid = this._belong2[target]
    if (bid1 === bid2) {
      this._updateFragment(start, end, 1)
      let cand = inclusive ? target : target + 1
      const vEnd = this._blockEnd2[vid]
      for (; cand < vEnd && !this._fragmentCounter2[cand]; cand++) {}
      if (cand < vEnd) {
        this._updateFragment(start, end, -1)
        return this._idToValue[cand]
      }

      let candVid = vid + 1
      for (; candVid < this._blockCount2 && !this._fragmentCounter1[candVid]; candVid++) {}
      if (candVid === this._blockCount2) {
        this._updateFragment(start, end, -1)
        return undefined
      }
      for (let j = this._blockStart2[candVid]; j < this._blockEnd2[candVid]; j++) {
        if (this._fragmentCounter2[j]) {
          this._updateFragment(start, end, -1)
          return this._idToValue[j]
        }
      }
    } else {
      this._updateFragment(start, this._blockEnd1[bid1], 1)
      this._updateFragment(this._blockStart1[bid2], end, 1)
      let cand = inclusive ? target : target + 1
      const vEnd = this._blockEnd2[vid]
      for (
        ;
        cand < vEnd &&
        !(this._fragmentCounter2[cand] + this._preSum2[bid2 * this._sortedLength + cand] - this._preSum2[(bid1 + 1) * this._sortedLength + cand]);
        cand++
      ) {}
      if (cand < vEnd) {
        this._updateFragment(start, this._blockEnd1[bid1], -1)
        this._updateFragment(this._blockStart1[bid2], end, -1)
        return this._idToValue[cand]
      }

      let candVid = vid + 1
      for (
        ;
        candVid < this._blockCount2 &&
        !(
          this._fragmentCounter1[candVid] +
          this._preSum1[bid2 * this._blockCount2 + candVid] -
          this._preSum1[(bid1 + 1) * this._blockCount2 + candVid]
        );
        candVid++
      ) {}
      if (candVid === this._blockCount2) {
        this._updateFragment(start, this._blockEnd1[bid1], -1)
        this._updateFragment(this._blockStart1[bid2], end, -1)
        return undefined
      }
      for (let j = this._blockStart2[candVid]; j < this._blockEnd2[candVid]; j++) {
        if (this._fragmentCounter2[j] + this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]) {
          this._updateFragment(start, this._blockEnd1[bid1], -1)
          this._updateFragment(this._blockStart1[bid2], end, -1)
          return this._idToValue[j]
        }
      }
    }

    throw new Error('unreachable')
  }

  private _countPrev(start: number, end: number, target: number, inclusive: boolean): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    target = this._getId(target)
    if (!inclusive) target--
    if (target < 0) return 0
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] <= target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd1[bid1]; i++) res += +(this._nums[i] <= target)
    for (let i = this._blockStart1[bid2]; i < end; i++) res += +(this._nums[i] <= target)
    const vid = this._belong2[target]
    for (let j = 0; j < vid; j++) {
      res += this._preSum1[bid2 * this._blockCount2 + j] - this._preSum1[(bid1 + 1) * this._blockCount2 + j]
    }
    for (let j = this._blockStart2[vid]; j <= target; j++) {
      res += this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
    }
    return res
  }

  private _countNext(start: number, end: number, target: number, inclusive: boolean): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    target = this._getId(target)
    if (!inclusive) target++
    if (target >= this._sortedLength) return 0
    const bid1 = this._belong1[start]
    const bid2 = this._belong1[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] >= target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd1[bid1]; i++) res += +(this._nums[i] >= target)
    for (let i = this._blockStart1[bid2]; i < end; i++) res += +(this._nums[i] >= target)
    const vid = this._belong2[target]
    for (let j = vid + 1; j < this._blockCount2; j++) {
      res += this._preSum1[bid2 * this._blockCount2 + j] - this._preSum1[(bid1 + 1) * this._blockCount2 + j]
    }
    for (let j = target; j < this._blockEnd2[vid]; j++) {
      res += this._preSum2[bid2 * this._sortedLength + j] - this._preSum2[(bid1 + 1) * this._sortedLength + j]
    }
    return res
  }
}

export { WaveletMatrixLikeOfflineDynamic }

if (require.main === module) {
  const ps = new WaveletMatrixLikeOfflineDynamic([3, 2, 1, 4, 7], [1, 2, 3, 4, 6, 7], 3)
  console.log(ps.kth(0, 4, 2))
  // ps.set(1, 6)
  console.log(ps.kth(1, 5, 2))
  console.log(ps.countRange(0, 5, 1, 7))
  // console.log(ps.count(0, 5, 2))

  test()
  performance()

  function performance(): void {
    const arr = Array.from({ length: 1e5 }, () => ~~(Math.random() * 1e9))
    const toadd = Array.from({ length: 1e5 }, () => ~~(Math.random() * 1e9))
    const allowed = new Set([...arr, ...toadd])
    console.time('init')
    const ps = new WaveletMatrixLikeOfflineDynamic(arr, allowed)
    for (let i = 0; i < arr.length; i++) {
      ps.set(0, toadd[i])
      ps.kth(0, arr.length, arr.length - 1)
    }
    console.timeEnd('init')
  }

  function test(): void {
    class Mocker implements WaveletMatrixLikeOfflineDynamic {
      private readonly nums: number[]

      constructor(allNums: number[], blockSize = 1000) {
        this.nums = allNums.slice()
      }

      set(pos: number, target: number): void {
        this.nums[pos] = target
      }

      kth(start: number, end: number, k: number): number | undefined {
        const tmp = this.nums.slice(start, end).sort((a, b) => a - b)
        return tmp[k]
      }

      lower(start: number, end: number, target: number): number | undefined {
        let res = -Infinity
        for (let i = start; i < end; i++) {
          if (this.nums[i] < target) res = Math.max(res, this.nums[i])
        }
        return res === -Infinity ? undefined : res
      }

      higher(start: number, end: number, target: number): number | undefined {
        let res = Infinity
        for (let i = start; i < end; i++) {
          if (this.nums[i] > target) res = Math.min(res, this.nums[i])
        }
        return res === Infinity ? undefined : res
      }

      floor(start: number, end: number, target: number): number | undefined {
        let res = -Infinity
        for (let i = start; i < end; i++) {
          if (this.nums[i] <= target) res = Math.max(res, this.nums[i])
        }
        return res === -Infinity ? undefined : res
      }

      ceiling(start: number, end: number, target: number): number | undefined {
        let res = Infinity
        for (let i = start; i < end; i++) {
          if (this.nums[i] >= target) res = Math.min(res, this.nums[i])
        }
        return res === Infinity ? undefined : res
      }

      count(start: number, end: number, target: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] === target) res++
        }
        return res
      }

      countRange(start: number, end: number, lower: number, upper: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] >= lower && this.nums[i] < upper) res++
        }
        return res
      }

      countLower(start: number, end: number, target: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] < target) res++
        }
        return res
      }

      countHigher(start: number, end: number, target: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] > target) res++
        }
        return res
      }

      countFloor(start: number, end: number, target: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] <= target) res++
        }
        return res
      }

      countCeiling(start: number, end: number, target: number): number {
        let res = 0
        for (let i = start; i < end; i++) {
          if (this.nums[i] >= target) res++
        }
        return res
      }

      get(index: number): number | undefined {
        return this.nums[index]
      }
    }

    const allNums = Array.from({ length: 5000 }, () => ~~(Math.random() * 1e9) - 8e5)
    const blockSize = ~~Math.sqrt(allNums.length) + 1
    const mocker = new Mocker(allNums.slice(), blockSize)
    const pointSetRangeKth = new WaveletMatrixLikeOfflineDynamic(allNums.slice(), allNums.slice())
    for (let i = 0; i < 1e4; i++) {
      const start = ~~(Math.random() * allNums.length)
      const end = ~~(Math.random() * (allNums.length - start)) + start
      const k = ~~(Math.random() * (end - start))
      const target = allNums[~~(Math.random() * allNums.length)]
      const lower = allNums[~~(Math.random() * allNums.length)]
      const upper = allNums[~~(Math.random() * allNums.length)]

      const res1 = mocker.kth(start, end, k)
      const res2 = pointSetRangeKth.kth(start, end, k)
      if (res1 !== res2) {
        console.error(`[${start}, ${end}) k=${k} res1=${res1} res2=${res2}`)
        throw new Error()
      }

      const res3 = mocker.lower(start, end, target)
      const res4 = pointSetRangeKth.lower(start, end, target)
      if (res3 !== res4) {
        console.error(`[${start}, ${end}) target=${target} res3=${res3} res4=${res4}`)
        throw new Error()
      }

      const res5 = mocker.higher(start, end, target)
      const res6 = pointSetRangeKth.higher(start, end, target)
      if (res5 !== res6) {
        console.error(`[${start}, ${end}) target=${target} res5=${res5} res6=${res6}`)
        throw new Error()
      }

      const resFloor = mocker.floor(start, end, target)
      const resFloor2 = pointSetRangeKth.floor(start, end, target)
      if (resFloor !== resFloor2) {
        console.error(`[${start}, ${end}) target=${target} resFloor=${resFloor} resFloor2=${resFloor2}`)
        throw new Error()
      }

      const resCeiling = mocker.ceiling(start, end, target)
      const resCeiling2 = pointSetRangeKth.ceiling(start, end, target)
      if (resCeiling !== resCeiling2) {
        console.error(`[${start}, ${end}) target=${target} resCeiling=${resCeiling} resCeiling2=${resCeiling2}`)
        throw new Error()
      }

      const res7 = mocker.countLower(start, end, target)
      const res8 = pointSetRangeKth.countLower(start, end, target)
      if (res7 !== res8) {
        console.error(`[${start}, ${end}) target=${target} res7=${res7} res8=${res8}`)
        throw new Error()
      }

      const res9 = mocker.count(start, end, target)
      const res10 = pointSetRangeKth.count(start, end, target)
      if (res9 !== res10) {
        console.error(`[${start}, ${end}) target=${target} res9=${res9} res10=${res10}`)
        throw new Error()
      }

      const res11 = mocker.countHigher(start, end, target)
      const res12 = pointSetRangeKth.countHigher(start, end, target)
      if (res11 !== res12) {
        console.error(`[${start}, ${end}) target=${target} res11=${res11} res12=${res12}`)
        throw new Error()
      }

      const res13 = mocker.countFloor(start, end, target)
      const res14 = pointSetRangeKth.countFloor(start, end, target)
      if (res13 !== res14) {
        console.error(`[${start}, ${end}) target=${target} res13=${res13} res14=${res14}`)
        throw new Error()
      }

      const res15 = mocker.countCeiling(start, end, target)
      const res16 = pointSetRangeKth.countCeiling(start, end, target)
      if (res15 !== res16) {
        console.error(`[${start}, ${end}) target=${target} res15=${res15} res16=${res16}`)
        throw new Error()
      }

      const res17 = mocker.countRange(start, end, lower, upper)
      const res18 = pointSetRangeKth.countRange(start, end, lower, upper)
      if (res17 !== res18) {
        console.error(`[${start}, ${end}) target=${target} res17=${res17} res18=${res18}`)
        throw new Error()
      }

      const resGet = mocker.get(start)
      const resGet2 = pointSetRangeKth.get(start)
      if (resGet !== resGet2) {
        console.error(`[${start}, ${end}) target=${target} resGet=${resGet} resGet2=${resGet2}`)
        throw new Error()
      }

      mocker.set(start, target)
      pointSetRangeKth.set(start, target)
    }

    console.log('pass')
  }
}
