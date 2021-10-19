const createElement = (str: string) => {
  const el = document.createElement('div')
  el.innerHTML = str
  return el.firstElementChild
}

const el = createElement(
  `<div class="container">
    <p>Hello!</p>
  </div>`
)

console.log(el!.className) // 'container'

export {}
