/**
 *
 * @param str  源字符串
 * @param target  关心的字符
 * @param k 可替换k次
 * @returns target最大连续长度
 */
const fix = (str: string, target: string, k: number): number => {
  let l = 0
  let r = 0
  let res = 0

  while (r < str.length) {
    if (str[r] !== target) k--

    while (k < 0) {
      l++
      if (str[l - 1] !== target) k++
    }

    res = Math.max(res, r - l + 1)
    r++
  }

  return res
}

export { fix }
