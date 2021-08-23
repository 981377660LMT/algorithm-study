// 序列化是指将一个数据结构转化为位序列的过程，因此可以将其存储在文件中或内存缓冲区中，
// 以便稍后在相同或不同的计算机环境中恢复结构。

class NAryNode<V> {
  val: V
  children?: NAryNode<V>[]

  constructor(val: V, children?: NAryNode<V>[]) {
    this.val = val
    this.children = children
  }
}

const tree: NAryNode<number> = {
  val: 1,
  children: [
    {
      val: 3,
      children: [
        { val: 5, children: undefined },
        { val: 6, children: undefined },
      ],
    },
    { val: 2, children: undefined },
    { val: 4, children: undefined },
  ],
}

// 思路一：dfs(适合一维数组转树形数组)
// 1,3,3,2,5,0,6,0,2,0,4,0
const serialize1 = (root: NAryNode<number>): string => {
  const res: string[] = []

  const dfs = (root: NAryNode<number> | null) => {
    if (!root) return
    res.push(root.val.toString())
    if (root.children) {
      res.push(root.children.length.toString())
      root.children.forEach(dfs)
    } else {
      res.push('0')
    }
  }
  dfs(root)

  return res.join(',')
}

// 思路：对root的每一个孩子， root.children.push(dfs())
const deserialize1 = (s: string): NAryNode<string> | null => {
  if (!s.length) return null
  function* genVal() {
    yield* s.split(',')
  }
  const gen = genVal()

  const dfs = () => {
    const val = gen.next().value!
    const childNum = Number(gen.next().value!)
    const root = new NAryNode<string>(val, childNum === 0 ? undefined : [])
    for (let i = 0; i < childNum; i++) {
      root.children?.push(dfs())
    }
    return root
  }
  return dfs()
}

// 思路二：bfs层序遍历(例如二叉树的反序列化根据层序遍历拼起来的)
const serialize2 = (root: NAryNode<number>): string => {
  if (!root) return ''
  const res: string[] = []
  const queue: NAryNode<number>[] = [root]
  while (queue.length) {
    const head = queue.pop()!
    res.push(head.val.toString())
    if (head.children) {
      res.push(head.children.length.toString())
      head.children.forEach(child => queue.push(child))
    } else {
      res.push('0')
    }
  }

  return res.join(',')
}

// 这也是反序列二叉树的解法:
// 1.建立childnode
// 2.push childnode到children
// 3.push childnode到queue
const deserialize2 = (s: string): NAryNode<string> | null => {
  if (!s.length) return null
  function* genVal() {
    yield* s.split(',')
  }

  const gen = genVal()
  const root = new NAryNode(gen.next().value!, [])
  // 记录节点和孩子的个数
  const queue: [NAryNode<string>, number][] = [[root, Number(gen.next().value!)]]
  while (queue.length) {
    const [head, num] = queue.shift()!
    for (let i = 0; i < num; i++) {
      const val = gen.next().value!
      const childNum = gen.next().value!
      const childNode = new NAryNode<string>(val, childNum === '0' ? undefined : [])
      head.children?.push(childNode)
      queue.push([childNode, Number(childNum)])
    }
  }

  return root
}

console.log(serialize1(tree))
console.log(serialize2(tree))
console.dir(deserialize1('1,3,3,2,5,0,6,0,2,0,4,0'), { depth: null })
// console.dir(deserialize2('1,3,4,0,2,0,3,2,6,0,5,0'), { depth: null })
