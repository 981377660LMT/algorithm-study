/**
 * @param raw  源字符串
 * @param need  关心的字符
 * @param k 可替换k次
 * @returns target最大连续长度
 */
function fix<A extends ArrayLike<T>, T = unknown>(raw: A, need: T, k: number): number {
  let left = 0
  let res = 0

  for (let right = 0; right < raw.length; right++) {
    if (raw[right] !== need) k--

    while (k < 0) {
      if (raw[left] !== need) k++
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}

export { fix }
