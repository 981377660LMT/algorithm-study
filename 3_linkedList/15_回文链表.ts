class Node {
  value: number
  next?: Node
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(2)
const e = new Node(1)
a.next = b
b.next = c
c.next = d
d.next = e

// 想一下回文字符串判断(双指针，一个前移，一个后移，判断相等)
// 解法1：遍历链表存val到数组，再用数组的方法判断
// 解法2：O(n) 时间复杂度和 O(1) 空间复杂度解决
// 两个字符串一个头插一个尾插

const isPalindrome = (head: Node) => {
  let str1 = ''
  let str2 = ''
  while (head) {
    str1 = str1 + head.value
    str2 = head.value + str2
    head = head.next!
  }

  return str1 === str2
}

console.dir(isPalindrome(a), { depth: null })

export {}
