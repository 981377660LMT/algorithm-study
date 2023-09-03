interface ITree<NoneLeafValue, Leaf extends object> {
  value: NoneLeafValue
  children: (ITree<NoneLeafValue, Leaf> | Leaf)[]
}

class TreeManager<NoneLeafValue, Leaf> {
  constructor(tree: ITree<NoneLeafValue, Leaf>) {}
}

export { TreeManager }
