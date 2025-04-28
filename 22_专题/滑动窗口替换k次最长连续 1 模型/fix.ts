/* eslint-disable no-param-reassign */

/**
 * 允许将数组中的任意字符替换为target字符k次，求target字符的最大连续长度.
 *
 * @param arr  源字符串
 * @param target  关心的字符
 * @param k 可替换k次
 * @returns target 最大连续长度
 */
function fix<A extends ArrayLike<T>, T = unknown>(arr: A, target: T, k: number): number {
  let left = 0
  let res = 0
  for (let right = 0; right < arr.length; right++) {
    k -= +(arr[right] !== target)
    while (k < 0) {
      k += +(arr[left] !== target)
      left++
    }
    res = Math.max(res, right - left + 1)
  }
  return res
}

export { fix }
