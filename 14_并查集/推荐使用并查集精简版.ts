// 并查集简化版：存储数组元素下标
// 并查集的root 不一定要是一个map 也可以用数组当map
// 初始化时自己指向自己
const useUnionFind = (size: number) => {
  const parent = Array.from<number, number>({ length: size }, (_, i) => i)

  const find = (val: number) => {
    while (parent[val] !== val) {
      val = parent[val]
    }
    return val
  }

  const union = (key1: number, key2: number) => {
    const root1 = find(key1)
    const root2 = find(key2)
    // 这一步优化很关键:总是让大的根指向小的根
    parent[Math.max(root1, root2)] = Math.min(root1, root2)
  }

  return { union, find }
}

export { useUnionFind }
