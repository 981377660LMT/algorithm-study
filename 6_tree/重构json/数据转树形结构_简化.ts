export {}

interface IDataItem {
  id: string
  /** 空字符串表示不存在父结点. */
  parentId: string
  name: string
}

interface ITreeNode extends IDataItem {
  children?: ITreeNode[]
}

function transform(arr: IDataItem[]): ITreeNode[] {
  const map = new Map<string, ITreeNode>()
  const roots: ITreeNode[] = []

  // 第一遍：建映射
  for (const item of arr) map.set(item.id, { ...item })

  // 第二遍：挂子节点
  for (const node of map.values()) {
    if (node.parentId === '') {
      roots.push(node)
    } else {
      const parent = map.get(node.parentId)!
      ;(parent.children ??= []).push(node)
    }
  }

  return roots
}

console.dir(
  transform([
    { id: 'A', name: 'A', parentId: '' },
    { id: 'A1', name: 'A1', parentId: 'A' },
    { id: 'A12', name: 'A12', parentId: 'A1' },
    { id: 'B', name: 'B', parentId: '' },
    { id: 'B2', name: 'B2', parentId: 'B' },
    { id: 'B22', name: 'B22', parentId: 'B2' }
  ]),
  { depth: null }
)
