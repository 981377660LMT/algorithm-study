import { BinaryTree } from '../Tree'
/**
 * Encodes a tree to a single string.
 *
 * @param {TreeNode} root
 * @return {string}
 * @description 容易层序遍历即可
 */
// const serialize = (root: TreeNode): string => {}

/**
 * Decodes your encoded data to tree.
 *
 * @param {(number | null)[]} data
 * @return {TreeNode}
 */
const deserializeNode = (data: (number | null)[]): BinaryTree | null => {
  if (!data.length) return null
  const genNode = (val?: number | null) => (val == null ? null : new BinaryTree(val))
  const root = new BinaryTree(data.shift()!)
  const queue: (BinaryTree | null)[] = [root]

  while (queue.length) {
    const head = queue.shift()
    if (head) {
      head.left = genNode(data.shift())
      head.right = genNode(data.shift())
      head.left && queue.push(head.left)
      head.right && queue.push(head.right)
    }
  }

  return root
}

// console.log(deserialize([1, 2, 3, null, null, 4, 5]))

export { deserializeNode }
