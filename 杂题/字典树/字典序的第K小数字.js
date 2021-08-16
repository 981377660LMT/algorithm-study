/**
 * @param {number} n
 * @param {number} k
 * @return {number}
 */
function findKthNumber(n, k) {
  /** 当前探索的出发点，在我们脑中的十叉树中，从这个数字开始往下先序遍历 */
  let result = 1

  /** 剩余探索步数，等这个值为 0 的时候，就说明 result 已经踩在了我们要的数字上 */
  let lastingSteps = k - 1 // 我们的 result 从 1 开始，就已经走了一步了
  while (lastingSteps > 0) {
    // 数一数先序遍历的话，从当前节点到自己的下一个兄弟节点，中间要经过几个子节点
    /** 最多能经过的子节点数，我们尝试扩大这个数，跳过尽量多的子节点 */
    let numberOfSubNode = 0
    /** 向下展开当前节点（得到 10 11 12 13 等子节点），得到的子树中左下角的那个节点 */
    let expandedChildrenStart = result
    /** 当前节点的下一个兄弟节点，比如 1 的，就是 2 */
    let expandedChildrenEnd = result + 1
    // 只要向下展开后，最左下角的那个子节点的数字的量级没有超过 n ，我们就尽量往下展开
    while (expandedChildrenStart <= n) {
      // 如果当前探索的出发点的子节点形成完全十叉树，Math.min 会取 expandedChildrenEnd，说明子树中就是有这么多个节点
      // Math.min 如果取到 n + 1 ，说明形成完全十叉树所需的节点比 n 还大，那么子树就是不完全的十叉树
      numberOfSubNode += Math.min(expandedChildrenEnd, n + 1) - expandedChildrenStart
      // 计算十叉树再展开一层子节点之后，子节点的数量
      // 如果这次展开之后，下一次 expandedChildrenEnd 大于 n + 1 ，就说明有一层节点数量不够形成完全十叉树了
      expandedChildrenStart *= 10
      expandedChildrenEnd *= 10
    }

    // 计算完当前节点到下一个兄弟节点之间有多少个，我们就能把这个量加到 result 上，从而把 result 对应的十叉树节点指针，移动到指向下一个兄弟节点
    // 如果发现展开后，当前位置的子节点的数量巨多，大于我们剩余的步数
    if (numberOfSubNode > lastingSteps) {
      // 就把当前探索出发点 result 这个指针在十叉树里下移一层
      result *= 10
      // 指针下移一层，等价于先序遍历走了一步，所以剩余步数减一
      lastingSteps -= 1
    } else {
      // 如果发现子节点数量小于剩余步数，那么我们就把指针往右移动，跳过整棵子树
      result += 1
      // 跳过子树相当于跳过这么多个子节点
      lastingSteps -= numberOfSubNode
    }
  }
  return result
}
