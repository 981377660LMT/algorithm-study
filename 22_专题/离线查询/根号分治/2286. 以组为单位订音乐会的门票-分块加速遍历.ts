// 2286. 以组为单位订音乐会的门票
// https://leetcode.cn/problems/booking-concert-tickets-in-groups/description/
//
// 对于 row 排 col 列的音乐厅设计一个音乐会买票系统，支持两种买票方式。
//
// 同一组的 k 位观众在同一排坐在一起。该方式对应 gather，称为聚集买票。
// 同一组的 k 位观众不一定坐在一起。该方式对应 scatter，称为分散买票。
//
// - 聚集买票要求在排数不超过 maxRow 的范围判断是否存在有 k 个空座位的排，
//   如果存在则需要计算排数最小的排的编号和这一排第一个空座位的编号，完成买票。
// - 分散买票要求在排数不超过 maxRow 的范围判断是否有至少 k 个空座位的票，
//   如果存在则需要计算买票的所有座位，完成买票。
//
// 1 <= row <= 5e4
// 1 <= col <= 1e9
// 0 <= maxRow <= n - 1
// gather 和 scatter 总 调用次数不超过 5e4 次。

// 如果维护的信息难以用线段树\平衡树维护，可以考虑分块.
// !暴力做法的瓶颈在于行可能很多，直接遍历复杂度太大，因此采用`分块加速遍历`.
// !每个块内维护sqrt(row)个行.
// 注意:1. 每次修改remain后需要重构块; 2. 将重构块提取成函数.

import { useBlock } from './SqrtDecomposition/useBlock'

class BookMyShow {
  private readonly _row: number
  private readonly _col: number
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array

  private readonly _remain: Uint32Array
  private readonly _blockSum: Float64Array
  private readonly _blockMax: Float64Array

  constructor(row: number, col: number) {
    const remain = new Uint32Array(row).fill(col)
    const { belong, blockStart, blockEnd, blockCount } = useBlock(remain)
    const blockMax = new Float64Array(blockCount)
    const blockSum = new Float64Array(blockCount)
    for (let i = 0; i < row; i++) {
      const bid = belong[i]
      blockMax[bid] = Math.max(blockMax[bid], remain[i])
      blockSum[bid] += remain[i]
    }

    this._row = row
    this._col = col
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd

    this._remain = remain
    this._blockSum = blockSum
    this._blockMax = blockMax
  }

  gather(k: number, maxRow: number): number[] {
    const start = 0
    const end = maxRow + 1
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]

    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        if (this._remain[i] >= k) {
          const startCol = this._col - this._remain[i]
          this._remain[i] -= k
          this._rebuildBlock(bid1)
          return [i, startCol]
        }
      }

      return []
    }

    for (let i = start; i < this._blockEnd[bid1]; i++) {
      if (this._remain[i] >= k) {
        const startCol = this._col - this._remain[i]
        this._remain[i] -= k
        this._rebuildBlock(bid1)
        return [i, startCol]
      }
    }

    for (let bid = bid1 + 1; bid < bid2; bid++) {
      if (this._blockMax[bid] < k) continue
      for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
        if (this._remain[i] >= k) {
          const startCol = this._col - this._remain[i]
          this._remain[i] -= k
          this._rebuildBlock(bid)
          return [i, startCol]
        }
      }
    }

    for (let i = this._blockStart[bid2]; i < end; i++) {
      if (this._remain[i] >= k) {
        const startCol = this._col - this._remain[i]
        this._remain[i] -= k
        this._rebuildBlock(bid2)
        return [i, startCol]
      }
    }

    return []
  }

  scatter(k: number, maxRow: number): boolean {
    const start = 0
    const end = maxRow + 1
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]

    if (bid1 === bid2) {
      let sum = 0
      for (let i = start; i < end; i++) sum += this._remain[i]
      if (sum < k) return false

      let rest = k
      for (let i = start; i < end; i++) {
        if (!rest) break
        const take = Math.min(this._remain[i], rest)
        this._remain[i] -= take
        rest -= take
      }
      this._rebuildBlock(bid1)
      return true
    }

    let sum = 0
    for (let i = start; i < this._blockEnd[bid1]; i++) sum += this._remain[i]
    for (let bid = bid1 + 1; bid < bid2; bid++) sum += this._blockSum[bid]
    for (let i = this._blockStart[bid2]; i < end; i++) sum += this._remain[i]
    if (sum < k) return false

    let rest = k
    for (let i = start; i < this._blockEnd[bid1]; i++) {
      if (!rest) break
      const take = Math.min(this._remain[i], rest)
      this._remain[i] -= take
      rest -= take
    }
    this._rebuildBlock(bid1)
    if (!rest) return true

    for (let bid = bid1 + 1; bid < bid2; bid++) {
      if (!this._blockSum[bid]) continue
      for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
        if (!rest) break
        const take = Math.min(this._remain[i], rest)
        this._remain[i] -= take
        rest -= take
      }
      this._rebuildBlock(bid)
      if (!rest) return true
    }

    for (let i = this._blockStart[bid2]; i < end; i++) {
      if (!rest) break
      const take = Math.min(this._remain[i], rest)
      this._remain[i] -= take
      rest -= take
    }
    this._rebuildBlock(bid2)
    return true
  }

  private _rebuildBlock(bid: number): void {
    this._blockMax[bid] = 0
    this._blockSum[bid] = 0
    for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
      const remain = this._remain[i]
      this._blockMax[bid] = Math.max(this._blockMax[bid], remain)
      this._blockSum[bid] += remain
    }
  }
}

/**
 * Your BookMyShow object will be instantiated and called as such:
 * var obj = new BookMyShow(n, m)
 * var param_1 = obj.gather(k,maxRow)
 * var param_2 = obj.scatter(k,maxRow)
 */
