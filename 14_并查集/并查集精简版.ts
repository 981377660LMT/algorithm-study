// 并查集简化版
// 为了简化将数组多加一位 第i位代表i
const UnionFind = (size: number) => {
  const parent = Array.from<number, number>({ length: size }, (_, i) => i)
  const find = (val: number) => {
    while (parent[val] !== val) {
      val = parent[val]
    }
    return val
  }
  const union = (pre: number, next: number) => {
    parent[next] = pre
  }
  return { union, find }
}
