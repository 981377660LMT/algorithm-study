// <body>
//   <ul>
//     <li>
//       <span>1</span>
//     </li>
//     <li>
//       <span>2</span>
//       <span>3</span>
//       <span id="span_3">4</span>
//     </li>
//   </ul>
// </body>
// 如果传入id = "span_3" 的元素，那么生成的xpath是body>ul[0]>li[1]>span[2]

// 反向查找
// 我们的参数是目标节点，我们的目标是冒泡到body，然后记录中间的节点即可。
function getXPath(node: HTMLElement) {
  const path: string[] = []
  helper(node, path)
  return path.reverse().join('>')

  function helper(node: HTMLElement, path: string[]): string {
    if (node === document.body) return `body${path}`

    const parentNode = node.parentNode as HTMLElement
    const index = Array.prototype.findIndex.call(parentNode.children, el => el === node)
    path.push(`${node.tagName.toLowerCase()}[${index}]`)
    return helper(parentNode, path)
  }
}
