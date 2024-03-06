/* eslint-disable arrow-body-style */
/* eslint-disable @typescript-eslint/ban-ts-comment */
/* eslint-disable max-len */
/* eslint-disable eqeqeq */

interface IOptions<NoneLeaf, Leaf> {
  isLeaf: (o: Leaf | NoneLeaf) => o is Leaf
  getChildren: (noneLeaf: NoneLeaf) => Array<Leaf | NoneLeaf>
}

class TreeManager<NoneLeaf, Leaf> {
  readonly root: NoneLeaf
  readonly isLeaf: (o: Leaf | NoneLeaf | undefined | null) => o is Leaf
  readonly getChildren: (node: Leaf | NoneLeaf | undefined | null) => Array<Leaf | NoneLeaf>

  constructor(root: NoneLeaf, options: IOptions<NoneLeaf, Leaf>) {
    this.root = root
    this.isLeaf = (o: Leaf | NoneLeaf | undefined | null): o is Leaf => o != undefined && options.isLeaf(o)
    this.getChildren = (node: Leaf | NoneLeaf | undefined | null) => (!node || options.isLeaf(node) ? [] : options.getChildren(node))
  }

  /**
   * 将结点筛选插入到 {@link parentPath} 对应结点的子结点位置之前.
   * @param pos 插入位置，默认为末尾.
   * @returns 返回插入结点的路径.
   */
  insertNode = (parentPath: number[], newNode: Leaf | NoneLeaf, pos?: number): number[] => {
    const parent = this.searchNode(parentPath)
    this._assertNoneLeaf(parent)
    const children = this.getChildren(parent)
    if (pos == undefined) pos = children.length
    children.splice(pos, 0, newNode)
    return [...parentPath, pos]
  }

  /**
   * 删除结点.不能删除根结点(即path不能为空).
   */
  removeNode = (nodePath: number[]): void => {
    if (!nodePath.length) throw new Error('path must not be empty when delete')
    const parentPath = nodePath.slice(0, -1)
    const parent = this.searchNode(parentPath)
    this._assertNoneLeaf(parent)
    const children = this.getChildren(parent)
    const childIndex = nodePath[nodePath.length - 1]
    children.splice(childIndex, 1)
  }

  /**
   * 更新整个结点.
   * 不能更新根结点(即path不能为空)，只能通过 {@link patchNode} 修改根结点.
   */
  updateNode = (nodePath: number[], newNode: Leaf | NoneLeaf): void => {
    if (!nodePath.length) throw new Error('path must not be empty when update')
    const parentPath = nodePath.slice(0, -1)
    const parent = this.searchNode(parentPath)
    this._assertNoneLeaf(parent)
    const children = this.getChildren(parent)
    const childIndex = nodePath[nodePath.length - 1]
    children.splice(childIndex, 1, newNode)
  }

  /**
   * 修改结点部分属性.
   */
  patchNode = (nodePath: number[], partialNode: Partial<Leaf> | Partial<NoneLeaf>): void => {
    const node = this.searchNode(nodePath)
    if (node == undefined) throw new Error('node not found')
    for (const key in partialNode) {
      if (Object.prototype.hasOwnProperty.call(partialNode, key)) {
        // @ts-ignore
        node[key] = partialNode[key]
      }
    }
  }

  searchNode = (nodePath: number[]): Leaf | NoneLeaf | undefined => {
    let res: Leaf | NoneLeaf = this.root
    for (let i = 0; i < nodePath.length; i++) {
      const childIndex = nodePath[i]
      const children = this.getChildren(res)
      if (childIndex < 0 || childIndex >= children.length) return undefined
      res = children[childIndex]
    }
    return res
  }

  /**
   * 对树进行剪枝，删除不需要的叶子节点.
   * @param shouldPrune 返回 true 则该叶子节点从树中删除.
   */
  pruneLeaf = (shouldPrune: (leaf: Leaf, path: number[]) => boolean): void => {
    const dfs = (cur: Leaf | NoneLeaf, parent: Leaf | NoneLeaf | undefined, path: number[]): void => {
      if (this.isLeaf(cur)) {
        this._assertNoneLeaf(parent)
        if (shouldPrune(cur, path)) {
          const children = this.getChildren(parent)
          children.splice(path[path.length - 1], 1)
        }
        return
      }

      const children = this.getChildren(cur)
      for (let i = children.length - 1; i >= 0; i--) {
        path.push(i)
        dfs(children[i], cur, path)
        path.pop()
      }
    }

    dfs(this.root, undefined, [])
  }

  /**
   * 对树进行剪枝，删除不需要的节点(无法删除根结点).
   * @param shouldPrune 返回 true 则该节点从树中删除.
   */
  pruneTree = (shouldPrune: (node: Leaf | NoneLeaf, path: number[]) => boolean): void => {
    const dfs = (cur: Leaf | NoneLeaf, parent: Leaf | NoneLeaf | undefined, path: number[]): void => {
      const children = this.getChildren(cur)
      for (let i = children.length - 1; i >= 0; i--) {
        path.push(i)
        dfs(children[i], cur, path)
        path.pop()
      }
      // 后序遍历时删除
      if (parent != undefined && shouldPrune(cur, path)) {
        const children = this.getChildren(parent)
        children.splice(path[path.length - 1], 1)
      }
    }

    dfs(this.root, undefined, [])
  }

  enumerateLeaf = (f: (leaf: Leaf) => void): void => {
    this.enumerateTree(node => {
      if (this.isLeaf(node)) f(node)
    })
  }

  enumerateTree = (f: (node: Leaf | NoneLeaf) => void): void => {
    const dfs = (cur: Leaf | NoneLeaf): void => {
      f(cur)
      if (this.isLeaf(cur)) return
      const children = this.getChildren(cur)
      for (let i = 0; i < children.length; i++) {
        dfs(children[i])
      }
    }

    dfs(this.root)
  }

  getPath = (node: Leaf | NoneLeaf): number[] => {
    const dfs = (cur: Leaf | NoneLeaf, path: number[]): boolean => {
      if (cur === node) return true
      const children = this.getChildren(cur)
      for (let i = 0; i < children.length; i++) {
        path.push(i)
        if (dfs(children[i], path)) return true
        path.pop()
      }
      return false
    }

    const path: number[] = []
    dfs(this.root, path)
    return path
  }

  getParent = (node: Leaf | NoneLeaf): NoneLeaf | undefined => {
    const path = this.getPath(node)
    if (!path.length) return undefined
    path.pop()
    return this.searchNode(path) as NoneLeaf
  }

  isNoneLeaf = (node: Leaf | NoneLeaf | undefined | null): node is NoneLeaf => {
    return node != undefined && !this.isLeaf(node)
  }

  print = (): void => {
    console.dir(this.root, { depth: null })
  }

  private _assertNoneLeaf(node: Leaf | NoneLeaf | undefined | null): asserts node is NoneLeaf {
    if (node == undefined || this.isLeaf(node)) {
      throw new TypeError(`node must be a none leaf node, but got ${node}`)
    }
  }
}

export { TreeManager }

if (require.main === module) {
  type NoneLeaf = { value: number; children: Array<Leaf | NoneLeaf> }
  type Leaf = { value: number }
  const root = {
    value: 0,
    children: [
      { value: 1 },
      { value: 2 },
      {
        value: 99,
        children: [{ value: 3 }, { value: 4 }, { value: 5 }, { value: 6 }]
      },
      { value: 7 },
      { value: 8 }
    ]
  } satisfies NoneLeaf

  const T = new TreeManager<NoneLeaf, Leaf>(root, {
    getChildren: node => node.children,
    isLeaf: (node): node is Leaf => !('children' in node)
  })

  T.print()
  T.pruneTree(node => T.isNoneLeaf(node))
  T.print()

  // 合并两颗树
  function testMergeFilterConfig(): void {
    const tree1 = {
      value: 0,
      children: [
        { value: 1 },
        { value: 2 },
        {
          value: 99,
          children: [{ value: 3 }, { value: 4 }, { value: 5 }, { value: 6 }]
        },
        { value: 7 },
        { value: 8 }
      ]
    } satisfies NoneLeaf

    const tree2 = {
      value: 0,
      children: [
        { value: 1 },
        { value: 2 },
        {
          value: 99,
          children: [{ value: 3 }, { value: 4 }, { value: 5 }, { value: 6 }]
        },
        { value: 7 },
        { value: 8 }
      ]
    } satisfies NoneLeaf

    const utils = new TreeManager<NoneLeaf, Leaf>(tree1, {
      getChildren: node => node.children,
      isLeaf: (node): node is Leaf => !('children' in node)
    })
    utils.insertNode([], tree2)

    utils.print()
  }

  testMergeFilterConfig()
}
