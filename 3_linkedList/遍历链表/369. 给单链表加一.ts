// 思路：找到最右边不为9的节点，加一；之后的所有9全部置为0
function plusOne(head: ListNode | null): ListNode | null {
  const dummy = new ListNode(0, head)
  let dummyP: ListNode | null = dummy
  let rightMostNoneNine = dummy

  while (dummyP) {
    if (dummyP.val !== 9) rightMostNoneNine = dummyP
    dummyP = dummyP.next
  }

  rightMostNoneNine.val += 1

  let noneNineP = rightMostNoneNine.next
  while (noneNineP) {
    noneNineP.val = 0
    noneNineP = noneNineP.next
  }

  return dummy.val === 1 ? dummy : dummy.next
}

export {}
