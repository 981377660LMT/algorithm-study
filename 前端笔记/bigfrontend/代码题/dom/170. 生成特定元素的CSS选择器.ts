/* eslint-disable prefer-destructuring */

// DOM选择器
// 输入一个 DOM， 返回这个 DOM 的选择器（selector）。
// 给定一个DOM结构，返回特定元素的CSS Selector。
// dfs记录路径

/**
 * 从当前结点向上遍历查找.
 */
function generateSelector(root: HTMLElement, target: HTMLElement): string {
  let node = target
  const path: string[] = []
  while (node !== root) {
    const tagName = node.tagName.toLowerCase()
    const id = node.id ? `#${node.id}` : ''
    const cur = `${tagName}${id}`
    path.push(cur)
    node = node.parentElement as HTMLElement
  }

  path.push(root.tagName.toLowerCase())
  return path.reverse().join(' > ')
}

// !这种方法需要遍历所有结点,不好
function _generateSelector(root: HTMLElement, target: HTMLElement): string {
  if (target.id) return `#${target.id}`
  let res = ''
  dfs(root, [])
  return res

  function dfs(cur: HTMLElement, path: string[]): void {
    if (cur === target) {
      res = path.join(' > ')
      return
    }

    const children = Array.from(cur.children)
    for (const child of children) {
      path.push(child.tagName.toLowerCase())
      dfs(child as HTMLElement, path)
      path.pop()
    }
  }
}

// 比如针对一下DOM结构，
// <div>
//   <p>BFE.dev</p>
//   <div>
//     is
//     <p>
//       <span>great. <button>click me!</button></span>
//     </p>
//   </div>
// </div>
// let selector = generateSelector(root, target) // 'button'
// expect(root.querySelector(selector)).toBe(target)

// selector = generateSelector(root, target) // 'div > div > p > button'
// expect(root.querySelector(selector)).toBe(target)
