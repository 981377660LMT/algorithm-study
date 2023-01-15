/* eslint-disable no-console */
// 简单双向绑定函数
function model(state: { value: string }, element: HTMLInputElement) {}

if (require.main === module) {
  const input = document.createElement('input')
  const state = { value: 'hello' }
  model(state, input)

  console.log(input.value) // hello
  input.value = 'world'
  console.log(state.value) // world
  state.value = 'hello world'
  input.dispatchEvent(new Event('change'))
  console.log(input.value) // hello world
}

export {}
