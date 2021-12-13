// 请你将 list1 中`下标`从 a 到 b 的全部节点都删除，并将list2 接在被删除节点的位置。
function mergeInBetween(
  list1: ListNode | null,
  a: number,
  b: number,
  list2: ListNode | null
): ListNode | null {
  const dummy = new ListNode(0, list1)
  let dummyP: ListNode = dummy

  for (let _ = 0; _ < a; _++) dummyP = dummyP.next!
  const end1 = dummyP

  for (let _ = 0; _ < b - a + 1; _++) dummyP = dummyP.next!
  const tmp = dummyP
  const end2 = tmp.next
  tmp.next = null

  end1.next = list2

  let headP = list2
  while (headP && headP.next) headP = headP.next!
  headP!.next = end2

  return dummy.next
}
