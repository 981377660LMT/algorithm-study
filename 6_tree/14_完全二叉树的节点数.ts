interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
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
    right: null,
  },
}

// const countNode = (root: TreeNode | null): number => {
//   if (root == null) return 0
//   return countNode(root.left) + countNode(root.right) + 1
// }
// 遍历树来统计节点是一种时间复杂度为 O(n) 的简单解决方案。
// 你可以设计一个更快的算法吗？
// 使用二分查找 复杂度O(log(n)^2)
// 思路:
// 1.注意到最后一层节点序号从2^h-1开始
const countNode = (root: TreeNode | null): number => {
  if (root == null) return 0
  const treeLevel = getTreeLevel(root)
  let l = 0
  let r = treeLevel ** 2 - 1
  while (l <= r) {
    const m = Math.floor((l + r) / 2)
    if (isExist(m, treeLevel, root)) {
      l = m + 1
    } else {
      r = m - 1
    }
  }

  return 2 ** treeLevel - 1 + l
}

// 快速求出树的高度(遍历链表)
// 第一行高为0
const getTreeLevel = (root: TreeNode) => {
  if (root == null) return 0
  let l = 0
  while (root.left) {
    root = root.left
    l++
  }
  return l
}

// console.log(getTreeLevel(bt))

// 二分法查找第level层序号为id的节点是否存在
// 用二进制表示树的路径 00 01 10 11  左走代表0 右走代表1
const isExist = (id: number, level: number, root: TreeNode) => {
  let l = 0
  let r = 2 ** level - 1
  while (l < r) {
    const m = Math.floor((l + r) / 2)
    if (id > m) {
      root = root.right!
      l = m + 1
    } else {
      root = root.left!
      r = m
    }
    console.log(root, l, r)
  }
  return !!root
}
console.log(isExist(0, 2, bt))
console.log(isExist(1, 2, bt))
console.log(isExist(2, 2, bt))
console.log(isExist(3, 2, bt))
console.log(isExist(4, 2, bt))

console.dir(countNode(bt), { depth: null })

export {}
