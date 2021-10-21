class Node {
  value: number
  next?: Node
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(4)
const b = new Node(2)
const c = new Node(1)
const d = new Node(3)
a.next = b
b.next = c
c.next = d

// 需要反转链表的一端节点
// 使用三节点

const sortList = (head: Node | undefined): Node | undefined => {
  const merge = (head1: Node | undefined, head2: Node | undefined): Node | undefined => {
    if (!head1) return head2
    if (!head2) return head1
    const dummy = new Node(-1)
    let dummyP: Node | undefined = dummy
    let head1P: Node | undefined = head1
    let head2P: Node | undefined = head2

    while (head1P && head2P) {
      if (head1P.value <= head2P.value) {
        dummyP.next = head1P
        head1P = head1P.next
      } else {
        dummyP.next = head2P
        head2P = head2P.next
      }
      dummyP = dummyP.next
    }

    dummyP.next = head1P || head2P
    return dummy.next
  }

  const mergeSort = (head: Node | undefined): Node | undefined => {
    if (!head || !head.next) return head

    // 快慢指针寻找中点
    let slow: Node | undefined = head
    let fast: Node | undefined = head
    while (fast && fast.next && fast.next.next) {
      fast = fast.next.next
      slow = slow!.next
    }

    const head2 = slow!.next
    slow!.next = undefined
    return merge(mergeSort(head), mergeSort(head2))
  }

  return mergeSort(head)
}

console.dir(sortList(a), { depth: null })
export {}
