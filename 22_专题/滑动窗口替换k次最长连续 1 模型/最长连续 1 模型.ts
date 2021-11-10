/**
 *
 * @param str  源字符串
 * @param target  关心的字符
 * @param k 可替换k次
 * @returns target最大连续长度
 */
const fix = (str: string, target: string, k: number): number => {
  let l = 0
  let res = 0

  for (let r = 0; r < str.length; r++) {
    if (str[r] !== target) k--

    while (k < 0) {
      if (str[l] !== target) k++
      l++
    }

    res = Math.max(res, r - l + 1)
  }

  return res
}

export { fix }
