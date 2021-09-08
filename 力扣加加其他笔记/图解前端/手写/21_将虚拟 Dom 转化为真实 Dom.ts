interface VirtualDom {
  tag: string
  attrs?: Record<string, string>
  children?: VirtualDom[]
}

const vDom = {
  tag: 'DIV',
  attrs: {
    id: 'app',
  },
  children: [
    {
      tag: 'SPAN',
      children: [{ tag: 'A', children: [] }],
    },
    {
      tag: 'SPAN',
      children: [
        { tag: 'A', children: [] },
        { tag: 'A', children: [] },
      ],
    },
  ],
}
// 把上诉虚拟Dom转化成下方真实Dom
// <div id="app">
//   <span>
//     <a></a>
//   </span>
//   <span>
//     <a></a>
//     <a></a>
//   </span>
// </div>

function dfs(vNode: VirtualDom | string) {
  if (typeof vNode === 'string') return document.createTextNode(vNode)
  const dom = document.createElement(vNode.tag)

  if (vNode.attrs) {
    for (const [key, value] of Object.entries(vNode)) {
      dom.setAttribute(key, value)
    }
  }

  vNode.children && vNode.children.forEach(child => dom.appendChild(dfs(child)))
  return dom
}
