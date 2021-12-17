// 临界点 定义为一个 局部极大值点 或 局部极小值点 。
// 注意：节点只有在同时存在前一个节点和后一个节点的情况下，才能成为一个 局部极大值点 / 极小值点 。
// 返回一个长度为 2 的数组 [minDistance, maxDistance] ，
// 其中 minDistance 是任意两个不同临界点之间的最小距离，
// maxDistance 是任意两个不同临界点之间的最大距离。如果临界点少于两个，则返回 [-1，-1] 。

// 1 <= Node.val <= 105

// 维护上一个和第一个临界点的位置
function nodesBetweenCriticalPoints(head: ListNode | null): number[] {
  let [minDist, maxDist] = [-1, -1]
  let [firstPeek, lastPeek] = [-1, -1]
  let index = 0
  let headP = head

  while (headP?.next?.next) {
    const [x, y, z] = [headP.val, headP.next.val, headP.next.next.val]

    // y是极值
    if (y > Math.max(x, z) || y < Math.min(x, z)) {
      if (lastPeek !== -1) {
        minDist = minDist === -1 ? index + 1 - lastPeek : Math.min(minDist, index + 1 - lastPeek)
        maxDist = Math.max(maxDist, index + 1 - firstPeek)
      }

      if (firstPeek === -1) {
        firstPeek = index + 1
      }

      lastPeek = index + 1
    }

    index++
    headP = headP.next
  }

  return [minDist, maxDist]
}
