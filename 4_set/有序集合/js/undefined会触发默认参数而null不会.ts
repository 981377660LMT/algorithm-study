class Node {
  val: number
  // left?: Node = undefined
  left: Node | undefined | null = null
  // left?: Node = undefined
  right: Node | undefined | null = null
  constructor(val = 0) {
    this.val = val
  }
}

const demo = new Node()

function* inorder(root: Node | null | undefined = demo): Generator<number> {
  if (!root) return
  yield* inorder(null)
  yield root.val
  yield* inorder(null)
}

for (const num of inorder()) {
  console.log(num)
}

export {}

// 只要传了undefined 就会取默认参数 陷入死循环
