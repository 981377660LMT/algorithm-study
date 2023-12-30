/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { PersistentArraySqrt } from '../../../22_专题/离线查询/根号分治/SqrtArray/PersistentArraySqrt'

// !完全可持久化数组:
// 每次修改都会生成一个新的数组

// 1 <= length <= 50000
// 题目最多进行50000 次set，snap，和 get的调用 。
// 0 <= index < length
// 0 <= snap_id < 我们调用 snap() 的总次数
// 0 <= val <= 10^9

// https://leetcode.cn/problems/snapshot-array/
class SnapshotArray {
  private readonly _git: PersistentArraySqrt<number>[] = []
  private _root: PersistentArraySqrt<number>

  constructor(length: number) {
    this._root = new PersistentArraySqrt(Array(length).fill(0))
  }

  set(index: number, val: number): void {
    this._root = this._root.set(index, val)
  }

  snap(): number {
    this._git.push(this._root)
    return this._git.length - 1
  }

  get(index: number, snapId: number): number {
    return this._git[snapId].get(index)!
  }
}
