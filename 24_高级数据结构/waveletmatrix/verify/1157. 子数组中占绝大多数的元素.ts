// 由于众数是绝对众数，即严格大于区间长度的一半。
// 因此区间内的中位数(偶数的话取左)一定是众数。
// 转换为求kth后求频率，即操作1操作2。

import { WaveletMatrix } from '../WaveletMatrix'
import { WaveletMatrixLikeOfflineDynamic } from '../WaveletMatrixLikeOfflineDynamic'

// https://leetcode.cn/problems/online-majority-element-in-subarray/description/
class MajorityChecker {
  private readonly _wm: WaveletMatrixLikeOfflineDynamic
  // private readonly _wm: WaveletMatrix

  constructor(arr: number[]) {
    this._wm = new WaveletMatrixLikeOfflineDynamic(arr, arr)
    // this._wm = new WaveletMatrix(new Uint32Array(arr))
  }

  query(left: number, right: number, threshold: number): number {
    const k = (right - left) >>> 1
    const kthMax = this._wm.kth(left, right + 1, k)!
    const freq = this._wm.count(left, right + 1, kthMax)
    return freq >= threshold ? kthMax : -1
  }
}
