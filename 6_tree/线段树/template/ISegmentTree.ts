// 接口仅供参考

interface ISegmentTreeNode<NodeValue = number> {
  left: number
  right: number
  value: NodeValue
  [key: string]: unknown
}

interface ISegmentTreeNodeWithLazy<NodeValue = number> extends ISegmentTreeNode<NodeValue> {
  isLazy: boolean
  lazyValue: NodeValue
}

interface ISegmentTree<TreeItem = number, QueryReturn = number> {
  update: (root: number, left: number, right: number, targetOrDelta: TreeItem) => void
  query: (root: number, left: number, right: number) => QueryReturn
}

export { ISegmentTree, ISegmentTreeNode, ISegmentTreeNodeWithLazy }
