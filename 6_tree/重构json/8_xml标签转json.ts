class TNode {
  name: string
  children?: TNode[]
  constructor(name: string, children?: TNode[]) {
    this.name = name
    this.children = children
  }
}

// 编写一个算法解析以下符号，转换为json树的结构
const str = `<xml><div><p><a/></p><p></p></div></xml><xml><div><p><a/></p><p></p></div></xml>`

// 是否可以转化为之前的题
const toTree = (str: string) => {
  // 虚拟根节点
  const root = new TNode('', [])

  const dfs = (str: string, parent: TNode) => {
    const regexp = /<(.*?)>(.*?)<\/\1>/g
    const match = [...str.matchAll(regexp)]

    if (match.length) {
      for (const group of match) {
        const name = group[1]
        const childStr = group[2]
        const root = new TNode(name, [])
        parent.children!.push(root)
        dfs(childStr, root)
      }
    } else {
      parent.children!.push(new TNode(str))
    }
  }
  dfs(str, root)

  return root.children
}
console.dir(toTree(str), { depth: null })
// console.log(toTree('<div><p><a/></p><p></p></div>'))
// console.log(toTree('<p><a/></p><p></p>'))
// console.log(toTree('<a/>'))
// console.log(str.match(match))
export {}
