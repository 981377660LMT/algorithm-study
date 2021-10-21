function model(state: { value: string }, element: HTMLInputElement) {
  element.value = state.value
  // type EventListenerOrEventListenerObject = EventListener | EventListenerObject;
  element.addEventListener('change', e => {
    state.value = e.target!.value
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

const input = document.createElement('input')
const state = { value: 'BFE' }
model(state, input)

console.log(input.value) // 'BFE'
state.value = 'dev'
console.log(input.value) // 'dev'
input.value = 'BFE.dev'
input.dispatchEvent(new Event('change'))
console.log(state.value) // 'BFE.dev'
