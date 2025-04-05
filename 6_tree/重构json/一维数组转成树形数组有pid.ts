interface IData {
  id: string
  pid: string
  name: string
}

interface ITreeNode extends IData {
  children?: ITreeNode[]
}

function trans(arr: IData[]): ITreeNode[] {
  const graph = new Map<string, IData[]>()
  for (const item of arr) {
    const { pid } = item
    if (!graph.has(pid)) graph.set(pid, [])
    graph.get(pid)!.push(item)
  }

  const roots: ITreeNode[] = []
  for (const item of arr) {
    if (item.pid === '') {
      const node: ITreeNode = { ...item }
      roots.push(node)
    }
  }
  roots.forEach(dfs)
  return roots

  function dfs(node: ITreeNode) {
    for (const child of graph.get(node.id) || []) {
      const childNode = { ...child }
      if (!node.children) node.children = []
      node.children.push(childNode)
      dfs(childNode)
    }
  }
}

const data: IData[] = [
  { id: '1', pid: '', name: 'A' },
  { id: '2', pid: '1', name: 'B' },
  { id: '3', pid: '1', name: 'C' },
  { id: '4', pid: '2', name: 'D' }
]

console.dir(trans(data), { depth: null })

export {}
