class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val: number = 0, left: TreeNode | null = null, right: TreeNode | null = null) {
    this.val = val
    this.left = left
    this.right = right
  }
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
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
}

const bt2: TreeNode = {
  val: 1,
  left: { val: 2, left: null, right: null },
  right: null,
}

const rightSideView = (root: TreeNode) => {
  const res: number[] = []
  const queue: TreeNode[] = [root]

  const bfs = (root: TreeNode | null, height: number) => {
    while (queue.length) {
      // 记录当前层级节点个数 对每层进行处理
      let length = queue.length

      while (length--) {
        const head = queue.shift()!
        //length长度为0的时候表明到了层级最后一个节点
        if (!length) {
          res.push(head.val)
        }
        head.left && queue.push(head.left)
        head.right && queue.push(head.right)
      }
    }
  }
  bfs(root, 0)

  return res
}

console.log(rightSideView(bt))
console.log(rightSideView(bt2))

export {}
