findClosestAnchor(document.querySelector('a > span')!) // a

function findClosestAnchor(node: Node) {
  let nodeP: Node | null = node
  while (nodeP) {
    if (nodeP.nodeName.toLowerCase() === 'a') return nodeP
    nodeP = nodeP.parentNode
  }
  return null
}
