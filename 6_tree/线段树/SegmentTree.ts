interface ISegmentTreeNode {
  left: number
  right: number
  [key: string]: any
}

interface ISegmentTree<TreeItem = number, QueryReturn = number> {
  update: (root: number, left: number, right: number, value: TreeItem) => void
  query: (root: number, left: number, right: number) => QueryReturn
}

export { ISegmentTree, ISegmentTreeNode }
