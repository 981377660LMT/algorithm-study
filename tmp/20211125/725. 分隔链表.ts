// 设计一个算法将链表分隔为 k 个连续的部分。
// 非常精巧的题，可以当均等分的模板记忆:每部分的长度应该尽可能的相等：任意两部分的长度差距不能超过 1 。
// 1.使用div mod 分割 序号在前mod的数可以多取1
// 2.游标headP定位 dummyP结点辅助查找
function splitListToParts(head: ListNode, k: number): ListNode[] {
  // if (!head) return []

  const res = Array<ListNode>(k)
  const size = getSize(head)
  const [div, mod] = [~~(size / k), size % k]

  let headP = head

  for (let i = 0; i < k; i++) {
    const dummy = new ListNode(0)
    let dummyP = dummy
    const steps = div + (i < mod ? 1 : 0)

    for (let j = 0; j < steps; j++) {
      dummyP.next = headP
      dummyP = dummyP.next
      headP = headP.next!
    }

    dummyP.next = null
    res[i] = dummy.next!
  }

  return res

  function getSize(head: ListNode | null): number {
    let size = 0
    let headP = head
    while (headP) {
      size++
      headP = headP.next
    }
    return size
  }
}

// 输入：head = [1,2,3,4,5,6,7,8,9,10], k = 3
// 输出：[[1,2,3,4],[5,6,7],[8,9,10]]
// 解释：
// 输入被分成了几个连续的部分，并且每部分的长度相差不超过 1 。前面部分的长度大于等于后面部分的长度。
export {}
