// 动态区间频率查询
// 区间加
// 查询区间某元素出现次数

import assert from 'assert'

import { SqrtDecomposition } from './SqrtDecomposition/SqrtDecomposition'

/**
 * 动态区间频率查询.
 */
class RangeFreqQueryDynamic {
  private readonly _sqrt: SqrtDecomposition<number, number, [v: number, same: boolean]>

  constructor(nums: number[]) {
    const n = nums.length

    this._sqrt = new SqrtDecomposition<number, number, [v: number, same: boolean]>(
      n,
      (_, left, right) => {
        const curNums = nums.slice(left, right)
        let sortedNums: number[] = []
        let lazyAdd = 0

        const created = () => {
          updated()
        }
        const updated = () => {
          sortedNums = [...curNums].sort((a, b) => a - b)
        }
        const updateAll = (lazy: number) => {
          lazyAdd += lazy
        }
        const updatePart = (left: number, right: number, lazy: number) => {
          for (let i = left; i < right; i++) {
            curNums[i] += lazy
          }
        }
        const queryAll = (queryArg: [v: number, same: boolean]) => {
          const [v, same] = queryArg
          if (same) {
            return bisectLeft(sortedNums, v - lazyAdd + 1) - bisectLeft(sortedNums, v - lazyAdd)
          }
          return sortedNums.length - bisectLeft(sortedNums, v - lazyAdd)
        }

        const queryPart = (left: number, right: number, queryArg: [v: number, same: boolean]) => {
          const [v, same] = queryArg
          if (same) {
            let res = 0
            for (let i = left; i < right; i++) {
              if (curNums[i] + lazyAdd === v) {
                res++
              }
            }
            return res
          }

          let res = 0
          for (let i = left; i < right; i++) {
            if (curNums[i] + lazyAdd >= v) {
              res++
            }
          }
          return res
        }

        return {
          created,
          updated,
          updateAll,
          updatePart,
          queryAll,
          queryPart
        }
      }
    )
  }

  /**
   * 区间`[left, right)`每个元素加上`value`.
   */
  update(left: number, right: number, value: number): void {
    this._sqrt.update(left, right, value)
  }

  /**
   * 查询区间`[left, right)`中元素`target`出现的次数.
   */
  rangeFreq(left: number, right: number, target: number): number {
    let res = 0
    this._sqrt.query(
      left,
      right,
      blockRes => {
        res += blockRes
      },
      [target, true]
    )
    return res
  }

  /**
   * 查询区间`[left, right)`中大于等于`floor`的元素出现的次数.
   */
  rangeFreqWithFloor(left: number, right: number, floor: number): number {
    let res = 0
    this._sqrt.query(
      left,
      right,
      blockRes => {
        res += blockRes
      },
      [floor, false]
    )
    return res
  }
}

function bisectLeft<T>(arr: ArrayLike<T>, value: T): number {
  let left = 0
  let right = arr.length - 1
  while (left <= right) {
    const mid = (left + right) >> 1
    if (arr[mid] < value) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return left
}

export { RangeFreqQueryDynamic }

if (require.main === module) {
  let rf = new RangeFreqQueryDynamic([1, 2, 2, 4, 5, 6, 7, 8, 9, 10])
  rf.update(0, 10, 1)
  assert.strictEqual(rf.rangeFreq(0, 10, 5), 1)
  rf.update(0, 10, 2)
  assert.strictEqual(rf.rangeFreq(0, 10, 5), 2)
  assert.strictEqual(rf.rangeFreqWithFloor(0, 10, 5), 9)

  const nums = Array.from({ length: 1e5 }, (_, i) => i)
  rf = new RangeFreqQueryDynamic(nums)
  console.time('time1')
  for (let i = 0; i < 1e5; i++) {
    rf.update(i, i + 1, i)
    rf.rangeFreq(0, i + 1, i)
  }
  console.timeEnd('time1') // 1.3s
}
