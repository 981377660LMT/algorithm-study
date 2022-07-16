import { BinaryTree } from '../分类/Tree'

/**
 * Encodes a tree to a single string.
 *
 * @param {TreeNode} root
 * @return {string}
 * @description 容易:层序遍历即可
 * 使用bfs,力扣就是这种输入 ['1,2,NULL']
 */
const serializeBFS = (root: BinaryTree | null): string => {
  if (!root) return 'NULL'

  const queue: (BinaryTree | null)[] = [root]
  const res: string[] = []

  while (queue.length) {
    const head = queue.shift()
    if (head != undefined) {
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
  const arr = data.split(',').map(val => (val === 'NULL' ? null : Number(val)))

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
const serializeDFS = (root: BinaryTree | null): string => {
  const res: string[] = []
  dfs(root)
  return res.join(',')

  function dfs(root: BinaryTree | null): void {
    if (!root) return
    res.push(root.val.toString())
    res.push(root.left ? '1' : '0')
    res.push(root.right ? '1' : '0')
    dfs(root.left)
    dfs(root.right)
  }
}

const deserializeDFS = (data: string): BinaryTree | null => {
  if (data === '') return null
  const dataGen = gen()
  return dfs()

  function dfs(): BinaryTree {
    const root = new BinaryTree(Number(next(dataGen)))
    const [hasLeft, hashRight] = [next(dataGen) === '1', next(dataGen) === '1']
    if (hasLeft) root.left = dfs()
    if (hashRight) root.right = dfs()
    return root
  }

  function* gen(): Generator<string, void, undefined> {
    yield* data.split(',')
  }

  function next<T>(iterator: Iterator<T>): T {
    return iterator.next().value
  }
}

// console.log(serializeDFS(deserializeDFS('1,2,NULL')))

/////////////////////////////////////////////////////////////////////////////////

/**
 * Decodes your encoded data to tree.
 *
 * @param {(number | null)[]} data
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

// function serialize(root: TreeNode | null): string {
//   return JSON.stringify(root)
// }

// /*
//  * Decodes your encoded data to tree.
//  */
// function deserialize(data: string): TreeNode | null {
//   return JSON.parse(data)
// }
