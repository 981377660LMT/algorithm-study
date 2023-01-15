/* eslint-disable prefer-destructuring */

// 给定一个DOM结构，返回特定元素的CSS Selector。
// dfs记录路径

function generateSelector(root: HTMLElement, target: HTMLElement): string {
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
