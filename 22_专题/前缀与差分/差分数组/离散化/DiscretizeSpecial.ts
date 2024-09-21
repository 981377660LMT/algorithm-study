/**
 * 给定元素 `0 ~ n-1`,对数组中的某些特殊元素进行离散化.
 *
 * @returns 返回离散化后的数组id和id对应的值.
 * 特殊元素的id为`0 ~ len(idToV)-1`, 非特殊元素的id为`-1`.
 */
function discretizeSpecial(
  n: number,
  isSpecial: (i: number) => boolean
): [vToId: number[], idToV: number[]] {
  const vToId = Array<number>(n)
  const idToV: number[] = []
  for (let i = 0; i < n; i++) {
    if (isSpecial(i)) {
      vToId[i] = idToV.length
      idToV.push(i)
    } else {
      vToId[i] = -1
    }
  }
  return [vToId, idToV]
}

export { discretizeSpecial }
