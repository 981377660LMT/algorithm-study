// 接口仅供参考

interface ISegmentTreeNode<NodeValue = number> {
  left: number
  right: number
  value: NodeValue
  [key: string]: any
}

interface ILazySegmentTreeNode<NodeValue = number> extends ISegmentTreeNode<NodeValue> {
  isLazy: boolean
  lazyValue: NodeValue
}

interface ISegmentTree<TreeItem = number, QueryReturn = number> {
  update: (root: number, left: number, right: number, targetOrDelta: TreeItem) => void
  query: (root: number, left: number, right: number) => QueryReturn
}

export { ISegmentTree, ISegmentTreeNode, ILazySegmentTreeNode }
