class Node {
  val: string
  children: Node[]
  constructor(val: string, children: Node[] = []) {
    this.val = val
    this.children = children
  }
}

class ThroneInheritance {
  private root: Node
  private deathSet: Set<string>
  private nameToNode: Map<string, Node> // 根据唯一名称直接找到对应节点

  // 创建这个多叉树的根节点
  constructor(kingName: string) {
    this.root = new Node(kingName)
    this.deathSet = new Set()
    this.nameToNode = new Map([[kingName, this.root]])
  }

  // 为某个特定的节点新增一个孩子
  birth(parentName: string, childName: string): void {
    const parent = this.nameToNode.get(parentName)
    const child = new Node(childName)
    parent?.children.push(child)
    this.nameToNode.set(childName, child)
  }

  // 在dfs遍历的序列中忽略该节点，其实也就是标记其死亡
  // 但注意不会影响树结构，其子孩子有权继位，优先级不变。
  death(name: string): void {
    this.deathSet.add(name)
  }

  // 在考虑death标志位的情况下，输出这个多叉树根节点的dfs遍历序列
  getInheritanceOrder(): string[] {
    const res: string[] = []
    this.preOrder(this.root, res, this.deathSet)
    return res
  }

  private preOrder(root: Node | undefined, res: string[], death: Set<string>) {
    if (!root) return
    if (!death.has(root.val)) res.push(root.val)
    for (const child of root.children) {
      this.preOrder(child, res, death)
    }
  }
}

export {}
