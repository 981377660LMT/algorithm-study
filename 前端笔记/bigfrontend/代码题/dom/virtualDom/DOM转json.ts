// 扩展思考:如果给定的不是一个 Dom 树结构 而是一段 html 字符串 该如何解析?

interface JSON {
  tag: string
  children: JSON[]
  [key: string]: any
}

function dom2Json(domtree: HTMLElement): JSON {
  const res: JSON = { tag: '', children: [] }
  res.tag = domtree.tagName
  res.children = []
  domtree.childNodes.forEach(child => res.children.push(dom2Json(child as HTMLElement)))
  return res
}

export {}
// <div>
//   <span>
//     <a></a>
//   </span>
//   <span>
//     <a></a>
//     <a></a>
//   </span>
// </div>

// 把上诉dom结构转成下面的JSON格式

// {
//   tag: 'DIV',
//   children: [
//     {
//       tag: 'SPAN',
//       children: [
//         { tag: 'A', children: [] }
//       ]
//     },
//     {
//       tag: 'SPAN',
//       children: [
//         { tag: 'A', children: [] },
//         { tag: 'A', children: [] }
//       ]
//     }
//   ]
// }
