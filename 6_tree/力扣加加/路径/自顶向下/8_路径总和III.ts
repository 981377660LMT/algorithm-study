// 深度优先
// 二叉树
interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
}

const bt: BinaryTree = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
}

// 是否存在一条子路径之和等于目标和
// 子路径 不需要从根节点开始，也不需要在叶子节点结束，但是路径方向必须是向下的（只能从父节点到子节点）。

// 思路:
// 前缀和就是到达当前元素的路径上，之前所有元素的和，即数列里的Sn。
// 如果在节点A和节点B处前缀和相差target，则位于节点A和节点B之间的元素之和是target
// 递归四步：
// 1.递归终止条件
// 2.当前层处理
// 3.进入下一层（可携带当前层变更后的状态进入下一层）
// 4.清理该层数据（回溯）
const pathSum = (root: BinaryTree | null, target: number): number => {
  if (!root) return 0
  let res = 0
  let curSum = 0
  // 这里设置初始值很关键 即target - curSum为0时加1
  const sumToCountMap = new Map<number, number>([[0, 1]])

  const dfs = (root: BinaryTree | null) => {
    if (!root) return

    curSum += root.val
    const matched = target - curSum
    const matchedCount = sumToCountMap.get(matched) || 0
    res += matchedCount
    sumToCountMap.set(curSum, sumToCountMap.get(curSum)! + 1 || 1)

    console.log(sumToCountMap, 66)
    root.left && dfs(root.left)
    root.right && dfs(root.right)

    // 回溯 清除记录值与总和值
    // 不回溯就相当于记录一路遍历下来的节点和
    sumToCountMap.set(curSum, sumToCountMap.get(curSum)! - 1)
    curSum -= root.val
  }
  dfs(root)

  return res
}

console.log(pathSum(bt, 7))
export {}
