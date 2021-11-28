// 反复删去链表中由 总和 值为 0 的连续节点组成的序列，直到不存在这样的序列为止。

// 就一个套路sum[i:j] = presum[j] - presum[i]
// 例如 2 3 -3 1 直接连2->1
// 总结：存储前缀和的`最右的一个结点`,直接一下全删完
function removeZeroSumSublists(head: ListNode | null): ListNode | null {
  const lastNodeOfPresum = new Map<number, ListNode>()
  const dummy = new ListNode(0, head)

  let curSum = 0
  let dummyP: ListNode | null = dummy
  while (dummyP) {
    curSum += dummyP.val
    lastNodeOfPresum.set(curSum, dummyP) // 开头是[0,-1]
    dummyP = dummyP.next
  }

  curSum = 0
  dummyP = dummy
  while (dummyP) {
    curSum += dummyP.val
    dummyP.next = lastNodeOfPresum.get(curSum)?.next ?? null
    dummyP = dummyP.next
  }

  return dummy.next
}

// 输入：head = [1,2,3,-3,-2]
// 输出：[1]
