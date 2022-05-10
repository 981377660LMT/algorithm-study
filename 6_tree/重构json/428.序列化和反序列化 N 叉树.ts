// 序列化是指将一个数据结构转化为位序列的过程，因此可以将其存储在文件中或内存缓冲区中，
// 以便稍后在相同或不同的计算机环境中恢复结构。

class MyNode {
  val: number
  children: MyNode[]
  constructor(val?: number) {
    this.val = val === undefined ? 0 : val
    this.children = []
  }
}

// dfs 序列化
// 1,3,3,2,5,0,6,0,2,0,4,0
const serialize = (root: MyNode): string => {
  const res: string[] = []
  dfs(root)
  return res.join(',')

  function dfs(root: MyNode | null): void {
    if (!root) return
    res.push(root.val.toString())
    res.push(root.children.length.toString())
    for (const child of root.children) {
      dfs(child)
    }
  }
}

// dfs 反序列化 对root的每一个孩子
// root.children.push(dfs())
const deserialize = (s: string): MyNode | null => {
  if (s.length === 0) return null
  const g = gen()
  return dfs()

  function dfs(): MyNode {
    const value = Number(next(g))
    const childCount = Number(next(g))
    const root = new MyNode(value)

    for (let _ = 0; _ < childCount; _++) {
      root.children.push(dfs())
    }

    return root
  }

  function next<T>(iter: Iterator<T>): T {
    return iter.next().value
  }

  function* gen(): Generator<string, void, undefined> {
    yield* s.split(',')
  }
}

// 总结
// 对于json(n叉树)
// dfs是从左往右深入处理
// bfs是对每层的键从上往下处理
