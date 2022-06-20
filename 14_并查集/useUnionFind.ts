/**
 *
 * @param iterable  元素类型为T的数组，用于初始化并查集
 * @description
 * 更加通用的并查集写法，调用add手动添加或union自动添加
 */
function useUnionFindMap<T = unknown>(iterable?: T[]) {
  let count = 0 // 连通分量个数
  const parent = new Map<T, T>()
  const rank = new Map<T, number>() // 各个帮派的大小
  for (const key of iterable ?? []) {
    add(key)
  }

  function add(key: T): boolean {
    if (parent.has(key)) return false
    parent.set(key, key)
    rank.set(key, 1)
    count++
    return true
  }

  // 如果key不在并查集，会自动add
  function find(key: T): T {
    if (!parent.has(key)) {
      add(key)
      return key
    }

    while (parent.has(key) && parent.get(key) !== key) {
      let p = parent.get(key)!
      // 进行路径压缩
      parent.set(p, parent.get(p)!)
      key = p
    }
    return key
  }

  function union(key1: T, key2: T): boolean {
    let root1 = find(key1)
    let root2 = find(key2)
    if (root1 === root2) return false
    if (rank.get(root1)! > rank.get(root2)!) {
      ;[root1, root2] = [root2, root1]
    }
    // rank优化:总是让大的根指向小的根
    parent.set(root1, root2)
    rank.set(root2, rank.get(root1)! + rank.get(root2)!)
    count--
    return true
  }

  function isConnected(key1: T, key2: T): boolean {
    return find(key1) === find(key2)
  }

  function getCount(): number {
    return count
  }

  function getParent(): Map<T, T> {
    return parent
  }

  function getRank(): Map<T, number> {
    return rank
  }

  function getRoots(): T[] {
    const res = new Set<T>()
    for (const key of parent.keys()) {
      const root = find(key)
      res.add(root)
    }
    return [...res]
  }

  return { add, union, find, isConnected, getCount, getRank, getRoots, getParent }
}

/**
 *
 * @param size 元素是0-size-1的并查集
 * @description
 * union不支持动态添加
 */
function useUnionFindArray(size: number) {
  let count = size
  const parents = Array.from<number, number>({ length: size }, (_, i) => i)
  const ranks = Array<number>(size).fill(1)

  function find(key: number) {
    while (parents[key] != undefined && parents[key] !== key) {
      // 进行路径压缩
      parents[key] = parents[parents[key]]
      key = parents[key]
    }
    return key
  }

  function union(key1: number, key2: number): boolean {
    let root1 = find(key1)
    let root2 = find(key2)
    if (root1 === root2) return false
    if (ranks[root1] > ranks[root2]) {
      ;[root1, root2] = [root2, root1]
    }
    // rank优化:总是让小的根指向大的根
    parents[root1] = root2
    ranks[root2] += ranks[root1]
    count--
    return true
  }

  function isConnected(key1: number, key2: number) {
    return find(key1) === find(key2)
  }

  function getCount(): number {
    return count
  }

  function getParents(): number[] {
    return parents
  }

  function getRanks(): number[] {
    return ranks
  }

  function getRoots(): number[] {
    const res = new Set<number>()
    for (let i = 0; i < size; i++) {
      const root = find(i)
      res.add(root)
    }
    return [...res]
  }

  return { union, find, isConnected, getCount, getRanks, getRoots, getParents }
}

if (require.main === module) {
  const uf = useUnionFindMap()
  uf.union(1, 2)
  console.log(uf.getCount())
}

export { useUnionFindArray, useUnionFindMap }
