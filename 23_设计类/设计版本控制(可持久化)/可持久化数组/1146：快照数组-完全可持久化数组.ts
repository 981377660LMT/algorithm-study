/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { PersistentArray, usePersistentArray } from './PersistentArray1'

// !完全可持久化数组:
// 每次修改都会生成一个新的数组

// 1 <= length <= 50000
// 题目最多进行50000 次set，snap，和 get的调用 。
// 0 <= index < length
// 0 <= snap_id < 我们调用 snap() 的总次数
// 0 <= val <= 10^9

const Q = 5e4 + 10

class SnapshotArray {
  private _snapId = 0
  private readonly _git = new Map<number, number>([[0, 0]]) // snapId => version
  private readonly _pArray: PersistentArray

  constructor(length: number) {
    this._pArray = usePersistentArray(length, Q)
  }

  /**
   * 会将指定索引 index 处的元素设置为 val
   */
  set(index: number, val: number): void {
    this._pArray.update(this._pArray.curVersion, index, val)
  }

  /**
   * 获取该数组的快照，并返回快照的编号 snap_id(快照号是调用 snap() 的总次数减去 1)。
   */
  snap(): number {
    this._git.set(this._snapId, this._pArray.curVersion)
    this._snapId++
    return this._snapId - 1
  }

  /**
   * 根据指定的 snap_id 调用 snap()，获取该数组在该快照下的指定 index 的值
   */
  get(index: number, snapId: number): number {
    const version = this._git.get(snapId)!
    return this._pArray.query(version, index)
  }
}

if (require.main === module) {
  const obj = new SnapshotArray(3)
  obj.set(0, 5)
  obj.snap()
  obj.set(0, 6)
  console.log(obj.get(0, 0))
}
