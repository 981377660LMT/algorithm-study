// 实现一个 render(template, context) 方法，将 template 中的占位符用 context 填充
const render = (template: string, context: Record<string, any>) => {
  return template.replace(/{{(\w+)}}/g, (_, g1, i) => {
    console.log(_, g1, i)
    return context[g1.trim()]
  })
}

console.log(render('{{name}}很厉害，才{{age}}岁', { name: 'bottle', age: '15' }))

export {}
