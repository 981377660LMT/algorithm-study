/**
 * @param {number} n
 * @param {number} k 注意：1 ≤ k ≤ n ≤ 10**9。
 * @return {number}
 * 求字典序第k个就是上图前序遍历访问的第k节点
 */
const findKthNumber = function (n: number, k: number): number {
  /**
   *
   * @param n
   * @param curNode
   * @param targetNode
   * @returns
   * 到右侧节点需要的步数
   */
  const countStepToRightNode = (n: number, curNode: number, targetNode: number) => {
    let step = 0
    while (curNode <= n) {
      // 比如n是195的情况195到100有96个数
      // targetNode-curNode 是一般情况
      // (n+1-curNode)表示正一行字典树没有满 到尽头了
      step += Math.min(targetNode, n + 1) - curNode
      curNode *= 10
      targetNode *= 10
    }
    return step
  }

  // 剩余移动步数
  let remainSteps = k
  remainSteps--
  let curNode = 1
  while (remainSteps) {
    // curNode + 1表示兄弟节点
    const step = countStepToRightNode(n, curNode, curNode + 1)

    if (step <= remainSteps) {
      // 到右边的需要的移动步数还小于等于k 还需要向右移动(从1到2,或者从10到11这种)
      remainSteps -= step
      curNode++
    } else {
      // 到右边的需要的移动步数已经大于k,需要向下移动(从1到10这种)
      remainSteps--
      curNode *= 10
    }
  }

  return curNode
}

console.log(findKthNumber(13, 2))
// 字典序的排列是 [1, 10, 11, 12, 13, 2, 3, 4, 5, 6, 7, 8, 9]，所以第二小的数字是 10。
