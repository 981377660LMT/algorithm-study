import { SqrtDecomposition } from '../SqrtDecomposition'

// https://leetcode.cn/problems/range-sum-query-mutable/
// Range Affine Range Sum

const INF = 2e15
class NumArray {
  private readonly _sqrt: SqrtDecomposition<number, number>

  constructor(nums: number[]) {
    this._sqrt = new SqrtDecomposition<number, number>(nums.length, (_, left, right) => {
      const curNums = nums.slice(left, right)
      let sum = 0
      let color = INF

      return {
        created() {
          this.updated()
        },
        updated() {
          sum = curNums.reduce((a, b) => a + b, 0)
        },
        updatePart(left, right, lazy) {
          for (let i = left; i < right; i++) {
            curNums[i] = lazy
          }
        },
        updateAll(lazy) {
          color = lazy
        },
        queryAll() {
          return color === INF ? sum : color * (right - left + 1)
        },
        queryPart(left, right) {
          let res = 0
          for (let i = left; i < right; i++) {
            res += color === INF ? curNums[i] : color
          }
          return res
        }
      }
    })
  }

  update(index: number, val: number): void {
    this._sqrt.update(index, index + 1, val)
  }

  sumRange(left: number, right: number): number {
    let res = 0
    this._sqrt.query(left, right + 1, blockRes => {
      res += blockRes
    })
    return res
  }
}

/**
 * Your NumArray object will be instantiated and called as such:
 * var obj = new NumArray(nums)
 * obj.update(index,val)
 * var param_2 = obj.sumRange(left,right)
 */
