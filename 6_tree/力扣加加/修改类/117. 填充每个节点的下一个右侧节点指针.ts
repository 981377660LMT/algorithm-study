import { BinaryTree } from '../Tree'

type AddPropRecursively<O extends object | null, P> = {
  [K in keyof O]: O[K] extends O | null ? AddPropRecursively<O[K], P> | null : O[K]
} &
  P

type NodeWithNext = AddPropRecursively<BinaryTree, { next: BinaryTree | null }>

const bt: NodeWithNext = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
      next: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
      next: null,
    },
    next: null,
  },
  right: {
    val: 3,
    left: null,
    right: {
      val: 7,
      left: null,
      right: null,
      next: null,
    },
    next: null,
  },
  next: null,
}

/**
 * @param {NodeWithNext} root
 * @return {NodeWithNext}
 * @description bfs记录上一个即可
 */
const connect = (root: NodeWithNext): NodeWithNext | null => {
  if (!root) return null
  const queue: [NodeWithNext, number][] = [[root, 0]]
  const res: number[][] = []

  while (queue.length) {
    const [head, level] = queue.shift()!
    head.next = queue[0] ? queue[0][0] : null
    if (!res[level]) {
      res[level] = [head.val]
    } else {
      res[level].push(head.val)
    }

    console.log(head.val, level)
    head.left && queue.push([head.left, level + 1])
    head.right && queue.push([head.right, level + 1])
  }

  return root
}

console.dir(connect(bt), { depth: null })
// 输出：[1,#,2,3,#,4,5,7,#]
export { AddPropRecursively }
