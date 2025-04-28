/**
 * 计算字符串 s 的最小周期。如果没有找到，返回 s 的长度。
 */
function minimalPeriod(s: string): number {
  const n = s.length
  const res = (s + s).indexOf(s, 1)
  return res !== -1 ? res : n
}

export { minimalPeriod }
