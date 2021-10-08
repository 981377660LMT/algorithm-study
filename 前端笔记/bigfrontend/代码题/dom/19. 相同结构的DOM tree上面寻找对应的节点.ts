/**
 * @param {HTMLElement} rootA
 * @param {HTMLElement} rootB - rootA and rootB are clone of each other
 * @param {HTMLElement} target
 * @description
 * 给定两个完全一样的DOM Tree A和B，以及A中的元素target，请找到B中对应的元素b。
 * 既然是DOM Tree，能否提供一个利用到DOM tree特性的解法？
   你的解法的时空复杂度是多少？
   这个问题可以出在一般的树结构上，DOM Tree只是一个特例。
 */
const findCorrespondingNode = (rootA: Element, rootB: Element, target: Element): Element | null => {
  // your code here
  if (rootA === target) return rootB
  for (let i = 0; i < rootA.children.length; i++) {
    const hit = findCorrespondingNode(rootA.children[i], rootB.children[i], target)
    if (hit) return hit
  }
  return null
}

const findCorrespondingNode2 = (
  rootA: Element,
  rootB: Element,
  target: Element
): Element | null => {
  // your code here
  if (rootA === target) return rootB

  const queueA = [rootA]
  const queueB = [rootB]

  while (queueA.length) {
    const curA = queueA.shift()!
    const curB = queueA.shift()!
    if (curA === target) return curB
    queueA.push(...Array.from(curA.children))
    queueB.push(...Array.from(curB.children))
  }

  return null
}

// createWalker API
const findCorrespondingNode3 = (rootA: Element, rootB: Element, target: Element): Node | null => {
  if (rootA === target) return rootB
  const treeWalkerA = document.createTreeWalker(rootA, NodeFilter.SHOW_ELEMENT)
  const treeWalkerB = document.createTreeWalker(rootB, NodeFilter.SHOW_ELEMENT)

  let curNodes: [Node | null, Node | null] = [treeWalkerA.currentNode, treeWalkerB.currentNode]
  while (curNodes[0] !== target) {
    curNodes = [treeWalkerA.nextNode(), treeWalkerB.nextNode()]
  }

  return curNodes[1]
}
