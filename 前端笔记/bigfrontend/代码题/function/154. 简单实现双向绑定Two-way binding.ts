/* eslint-disable no-param-reassign */
/* eslint-disable no-console */

// 简单双向绑定函数
function model(state: { value: string }, element: HTMLInputElement) {
  element.value = state.value // init

  // state -> DOM
  Object.defineProperty(state, 'value', {
    get() {
      return element.value
    },
    set(v) {
      element.value = v
    }
  })

  // DOM -> state
  element.addEventListener('change', e => {
    state.value = (e.target as HTMLInputElement).value
  })
}

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
