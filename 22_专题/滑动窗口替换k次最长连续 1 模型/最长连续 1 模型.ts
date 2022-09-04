/* eslint-disable no-param-reassign */

/**
 * @param arr  源字符串
 * @param need  关心的字符
 * @param k 可替换k次
 * @returns need 最大连续长度
 */
function fix<A extends ArrayLike<T>, T = unknown>(arr: A, need: T, k: number): number {
  let left = 0
  let res = 0

  for (let right = 0; right < arr.length; right++) {
    if (arr[right] !== need) k--

    while (k < 0) {
      if (arr[left] !== need) k++
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}

export { fix }
