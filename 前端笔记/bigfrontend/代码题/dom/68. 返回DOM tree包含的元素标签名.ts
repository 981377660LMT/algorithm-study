// 给定一个DOM tree，返回其中包含的所有的标签名。

/**
 * @param {HTMLElement} tree
 * @return {string[]}
 */
function getTags(tree: HTMLElement): string[] {
  const res = new Set<string>()

  const dfs = (node: Element) => {
    res.add(node.tagName.toLowerCase())
    for (const child of Array.from(node.children)) {
      dfs(child)
    }
  }
  dfs(tree)

  return [...res]
}

function getTags2(tree: HTMLElement): string[] {
  const res = new Set<string>([tree.nodeName])
  tree.querySelectorAll('*').forEach(node => res.add(node.nodeName))
  return [...res].map(s => s.toLowerCase())
}

function getTags3(tree: HTMLElement): string[] {
  const res = new Set<string>([tree.nodeName])
  const walker = document.createTreeWalker(tree, NodeFilter.SHOW_ELEMENT)

  let cur: Node | null = walker.currentNode
  while (cur) {
    res.add(cur.nodeName)
    cur = walker.nextNode()
  }

  return [...res].map(s => s.toLowerCase())
}

function getTags4(tree: HTMLElement): string[] {
  const nodeList = document.querySelectorAll('*')
  return [...new Set(Array.from(nodeList).map(node => node.tagName))]
}
