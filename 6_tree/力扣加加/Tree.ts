class BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
  constructor(val: number) {
    this.val = val
    this.left = null
    this.right = null
  }
}

const bt: BinaryTree = {
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

export { BinaryTree, bt }
