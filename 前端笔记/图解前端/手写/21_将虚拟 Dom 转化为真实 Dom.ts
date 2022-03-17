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

// leetcode`克隆图`
function dfs(vNode: VirtualDom | string): HTMLElement | Text {
  if (typeof vNode === 'string') return document.createTextNode(vNode)
  const res = document.createElement(vNode.tag)

  for (const [key, value] of Object.entries(vNode.attrs ?? {})) {
    res.setAttribute(key, value)
  }

  vNode.children?.forEach(child => res.appendChild(dfs(child)))
  return res
}

export {}
