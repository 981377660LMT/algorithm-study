import { bisectLeft } from '../../9_排序和搜索/二分/7_二分搜索寻找最左插入位置'
import { bisectRight } from '../../9_排序和搜索/二分/7_二分搜索寻找最插右入位置'

// 所有range都表示左闭右开
class RangeModule {
  private range: number[]

  constructor() {
    this.range = []
  }

  /**
   *
   * @param left
   * @param right
   * 添加区间需要先查一下会不会和已有的区间和交集，如果有则融合
   * [1,2,3,5,8,12]， 代入 left = 3, right = 5，此时需要保持不变，
   * 就不难知道应该用 bisect_left 还是 bisect_right
   * 添加时 l=2 r=4 都push
   */
  addRange(left: number, right: number): void {
    this.overlap(left, right, false)
  }

  /**
   *
   * @param left
   * @param right
   * [1,2,3,5,8,12]， 代入 left = 3, right = 5，此时需要保持不变，
   * 就不难知道应该用 bisect_left 还是 bisect_right
   * 删除时 l=2 r=4 都不push
   */
  removeRange(left: number, right: number): void {
    this.overlap(left, right, true)
  }

  /**
   *
   * @param left
   * @param right
   * @returns
   * 例如[1,20,25,30] 里询问[1,20)
   * l=1 r=1
   */
  queryRange(left: number, right: number): boolean {
    const l = bisectRight(this.range, left)
    const r = bisectLeft(this.range, right)
    // 是奇数并且在一个区间
    return (l & 1) === 1 && l == r
  }

  private overlap(left: number, right: number, idOdd: boolean): void {
    const l = bisectLeft(this.range, left)
    const r = bisectRight(this.range, right)
    const merge: number[] = []
    if ((l & 1) === Number(idOdd)) merge.push(left)
    if ((r & 1) === Number(idOdd)) merge.push(right)
    // 修改 ranges 的 [l:r-1] 部分
    this.range.splice(l, r - l, ...merge)
  }

  // [1,2,3,5,8,12]
  static main() {
    const rangeModule = new RangeModule()
    rangeModule.addRange(10, 20)
    rangeModule.removeRange(14, 16)
    console.log(rangeModule.queryRange(10, 14)) // true
    console.log(rangeModule.queryRange(13, 15)) // false
    console.log(rangeModule.queryRange(16, 17)) // true
  }
}

RangeModule.main()
