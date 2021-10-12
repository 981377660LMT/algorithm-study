// 请实现函数model(state, element)，使得state.value和HTMLInputElement element联动。
function model(state: { value: string }, element: HTMLInputElement) {
  // your code here
  element.value = state.value
  // type EventListenerOrEventListenerObject = EventListener | EventListenerObject;
  element.addEventListener('change', e => {
    // @ts-ignore
    state.value = e.target.value
  })

  Object.defineProperty(state, 'value', {
    get() {
      return element.value
    },
    set(v) {
      element.value = v
    },
  })
}

if (require.main === module) {
  const input = document.createElement('input')
  const state = { value: 'BFE' }
  model(state, input)

  console.log(input.value) // 'BFE'
  state.value = 'dev'
  console.log(input.value) // 'dev'
  input.value = 'BFE.dev'
  input.dispatchEvent(new Event('change'))
  console.log(state.value) // 'BFE.dev'
}
