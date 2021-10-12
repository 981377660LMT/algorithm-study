// Lit-HTML 允许您用 JavaScript 编写 HTML 模板，
// 然后有效地将这些模板与数据一起进行渲染和重新渲染，以创建和更新 DOM

function html(strs: TemplateStringsArray, ...keys: string[]): string {
  // 一个萝卜一个坑
  return strs.map((str, index) => `${str}${keys[index] ?? ''}`).join('')
}

// render the result from html() into the container
function render(result: string, container: HTMLElement) {
  // your code here
  container.innerHTML = result
}

if (require.main === module) {
  const helloTemplate = (name: string) => html`<div>Hello ${name}!</div>`
  console.log(helloTemplate('121'))
  // This renders <div>Hello Steve!</div> to the document body
  // render(helloTemplate('Steve'), document.body)

  // // This updates to <div>Hello Kevin!</div>, but only updates the ${name} part
  // // 魔法发生在第二次的render() ，其只更新了需要更新的部分。
  // render(helloTemplate('Kevin'), document.body)
}

// 在JavaScript中，??运算符被称为nullish coalescing运算符(零合并操作符)。
// 如果第一个参数不是null/undefined，这个运算符将返回第一个参数，否则，它将返回第二个参数
