import { BinaryTree } from '../../6_tree/力扣加加/Tree'
import { deserializeNode } from '../../6_tree/力扣加加/构建类/297.二叉树的序列化与反序列化'

// 任务调度优化是计算机性能优化的关键任务之一
// 任务之间是存在依赖关系的
// root 为根任务，root.left 和 root.right 为他的两个前导任务（可能为空
// 两个 CPU 核，即我们可以同时执行两个任务，但是同一个任务不能同时在两个核上执行。给定这颗任务树，请求出所有任务执行完毕的最小时间。

// 每个节点的任务执行时间的最小值，应该是Max(time(node.left),time(node.right),preTime/2) + node.val
// preTime是只有一个CPU时的前置任务总时间
// 每个节点的任务执行时间的最小值：左子树执行完成的最小时间、右子树执行完成的最小时间、左右子树全部节点并行执行的时间，三者的最大值，再加上当前节点的任务时间。
function minimalExecTime(root: BinaryTree | null): number {
  /**
   *
   * @param root
   * @param cpuCount
   * @returns 最小耗时；子树结点之和(串行时间)
   */
  function dfs(root: BinaryTree | null, cpuCount: number): [minTime: number, subtreeSum: number] {
    if (!root) return [0, 0]

    const [leftMin, leftSum] = dfs(root.left, cpuCount)
    const [rightMin, rightSum] = dfs(root.right, cpuCount)
    const sum = leftSum + rightSum
    // 注意子树的任务运行不能比sum / cpuCount还小(全部并行)
    const minTime = Math.max(leftMin, rightMin, sum / cpuCount) + root.val

    return [minTime, sum + root.val]
  }

  return dfs(root, 2)[0]
}

console.log(minimalExecTime(deserializeNode([15, 21, null, 24, null, 27, 26])))
// 87
console.log(minimalExecTime(deserializeNode([1, 3, 2, null, null, 4, 4])))
// 输出：7.5
// 1+（3+2+4+4）/2
