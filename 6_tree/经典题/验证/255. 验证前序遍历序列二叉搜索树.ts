// BST前序遍历：局部单减，整体单增
// 栈内局部递减，表示左子树递归；出现递增点，表示升级切换到右子树了；开始整体递增
// 如果检测到某次入栈元素小于上一次出栈元素，则说明不是二叉搜索树

// 上一次出栈元素是被右叶子结点弹出来的，表示增大阶段，后面的数不可能比它再小了
function verifyPreorder(preorder: number[]): boolean {
  const stack: number[] = []
  let prePop = -Infinity

  for (const cur of preorder) {
    if (prePop > cur) return false

    while (stack.length >= 1 && stack[stack.length - 1] < cur) {
      prePop = stack.pop()!
    }

    stack.push(cur)
  }

  return true
}

console.log(verifyPreorder([5, 2, 1, 3, 6]))
