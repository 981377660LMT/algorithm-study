import { BinaryTree } from '../Tree'

// https://leetcode-cn.com/problems/xu-lie-hua-er-cha-shu-lcof/solution/shou-hua-tu-jie-dfshe-bfsliang-chong-jie-fa-er-cha/
/**
 * Encodes a tree to a single string.
 *
 * @param {TreeNode} root
 * @return {string}
 * @description 容易:层序遍历即可
 * 使用bfs,力扣就是这种输入
 */
const serializeBFS = (root: BinaryTree | null): string => {
  if (!root) return 'NULL'

  const queue: (BinaryTree | null)[] = [root]
  const res: string[] = []

  while (queue.length) {
    const head = queue.shift()
    if (head) {
      res.push(root.val.toString())
      queue.push(head.left)
      queue.push(head.right)
    } else {
      res.push('NULL')
    }
  }

  return res.join(',')
}

const deserializeBFS = (data: string): BinaryTree | null => {
  if (data === 'NULL') return null
  const arr = data.split(',').map(v => (v === 'NULL' ? null : Number(v)))

  const inner = (arr: (number | null)[]): BinaryTree | null => {
    if (!arr.length) return null
    const genNode = (val?: number | null) => (val == null ? null : new BinaryTree(val))
    const root = new BinaryTree(arr.shift()!)
    const queue: (BinaryTree | null)[] = [root]

    while (queue.length) {
      const head = queue.shift()
      if (head) {
        head.left = genNode(arr.shift())
        head.right = genNode(arr.shift())
        head.left && queue.push(head.left)
        head.right && queue.push(head.right)
      }
    }

    return root
  }

  return inner(arr)
}
// console.log(serializeBFS(null))
// console.log(deserializeBFS('NULL'))

////////////////////////////////////////////////////////////////////////////

// 使用dfs
// 简洁的做法
const serializeDFS = (root: BinaryTree | null): string => {
  if (!root) return '_'
  const res: string[] = []
  const dfs = (root: BinaryTree | null) => {
    if (!root) return res.push('_')
    res.push(root.val.toString())
    dfs(root.left)
    dfs(root.right)
  }
  dfs(root)
  return res.join(',')
}

const deserializeDFS = (data: string) => {
  const arr = data.split(',').map(v => (v === '_' ? null : Number(v)))
  const dfs = (arr: (number | null)[]) => {
    const val = arr.shift()
    if (val == null) return null
    const node = new BinaryTree(val)
    node.left = dfs(arr)
    node.right = dfs(arr)
    return node
  }

  return dfs(arr)
}

// console.log(serializeDFS(deserializeDFS('1,2,NULL')))

/////////////////////////////////////////////////////////////////////////////////

/**
 * Decodes your encoded data to tree.
 *
 * @param {(number | null)[]} data
 * @return {TreeNode}
 * 使用bfs
 */
function deserializeNode(data: (number | null)[]): BinaryTree | null {
  if (!data.length) return null
  const genNode = (val?: number | null) => (val == null ? null : new BinaryTree(val))
  const root = new BinaryTree(data.shift()!)
  const queue: (BinaryTree | null)[] = [root]

  while (queue.length) {
    const cur = queue.shift()
    if (cur) {
      cur.left = genNode(data.shift())
      cur.right = genNode(data.shift())
      cur.left && queue.push(cur.left)
      cur.right && queue.push(cur.right)
    }
  }

  return root
}

// console.log(deserialize([1, 2, 3, null, null, 4, 5]))

export { deserializeNode }
