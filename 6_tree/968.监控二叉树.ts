import { BinaryTree } from './分类/Tree'
import { deserializeNode } from './重构json/297.二叉树的序列化与反序列化'

enum State {
  NeedShoot = 0,
  AlreadyShoot = 1,
  Monitor = 2,
}

/**
 * @param {BinaryTree} root
 * @return {number}
 * 节点上的每个摄影头都可以监视其父对象、自身及其直接子对象。计算监控树的所有节点所需的最小摄像头数量。
 */
const minCameraCover = (root: BinaryTree): number => {
  let res = 0
  const dfs = (root: BinaryTree | null): State => {
    // 空节点不需要被人拍也不用拍别人，直接返回被拍了就好
    if (!root) return State.AlreadyShoot
    const left = dfs(root.left)
    const right = dfs(root.right)

    // 我装个摄像机:左儿子或者右儿子需要被拍
    if (left === State.NeedShoot || right === State.NeedShoot) {
      res++
      return State.Monitor
    }

    // 我被拍了:左儿子或者右儿子装了摄像机，
    if (left === State.Monitor || right === State.Monitor) {
      return State.AlreadyShoot
    }

    // 我需要被拍:左儿子和右儿子都是被拍的，都没有摄像机
    return State.NeedShoot
  }

  // 根节点需不需要被拍
  if (dfs(root) === State.NeedShoot) res++

  // console.dir(root, { depth: null })
  return res
}

console.log(minCameraCover(deserializeNode([0, 0, null, 0, 0])!))
export default 1
