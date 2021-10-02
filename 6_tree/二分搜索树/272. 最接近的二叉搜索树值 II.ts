import { ArrayDeque } from '../../2_queue/Deque'
import { BinaryTree } from '../力扣加加/Tree'

/**
 *
 * @param root
 * @param target
 * @param k 请在该二叉搜索树中找到最接近目标值 target 的 k 个值。
 * 把Queue改成priority queue还可以解决不是二叉搜索树的情况，时间复杂度会变成 nlogk。 另外空间复杂度为O(h+k), h平均为log n.
 */
function closestKValues(root: BinaryTree | null, target: number, k: number): number[] {
  const queue = new ArrayDeque(10000)
  inorder(root)
  return [...queue]

  function inorder(root: BinaryTree | null) {
    if (!root) return

    inorder(root.left)

    if (queue.length < k) {
      queue.push(root.val)
    } else {
      if (Math.abs(root.val - target) < Math.abs(queue.front()! - target)) {
        queue.shift()
        queue.push(root.val)
      } else return
    }

    inorder(root.right)
  }
}

// 在中序遍历的同时，当list的个数小于k，直接添加。
// 等于的时候，就拿第一个与当前节点各自与目标值target的差值的绝对值进行比较。
// 第一个与目标值的差值的绝对大于当前节点与目标值的差值的绝对值，
// 则移除首位，然后添加当前节点，保持list的个数是k个。
// 反之则返回，无需再遍历后面的节点。
