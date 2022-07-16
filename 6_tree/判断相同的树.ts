import { BinaryTree } from './分类/Tree'
import { deserializeNode } from './重构json/297.二叉树的序列化与反序列化'

// 相同的树
function isSameTree(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
  if (root1 == null && root2 == null) return true
  if (root1 == null || root2 == null) return false // 一个空一个不空

  return (
    root1.val === root2.val &&
    isSameTree(root1.left, root2.left) &&
    isSameTree(root1.right, root2.right)
  )
}

if (require.main === module) {
  console.dir(
    isSameTree(deserializeNode([1, 2, 2, 3, 4, 4, 3]), deserializeNode([1, 2, 2, 3, 4, 4, 3])),
    { depth: null }
  )
}

export {}

// 同样的逻辑适用于对象的 deepEqual :
// 当前是值类型,判断相等;当前是引用类型,Object.keys遍历看子代相等(键/值都要等)
function isSameObject(o1: unknown, o2: unknown): boolean {
  if (Object.is(o1, o2)) return true

  if (isObject(o1) && isObject(o2)) {
    const keys1 = Object.keys(o1)
    const keys2 = Object.keys(o2)
    if (keys1.length !== keys2.length) return false

    for (const key of keys1) {
      if (!isSameObject(o1[key], o2[key])) return false
    }

    return true
  }

  return false
}

function isObject(o: unknown): o is Record<any, any> {
  return typeof o === 'object' && o !== null
}
