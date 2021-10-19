export {}
const fn = () => console.log('!')
document.body.addEventListener('click', fn)
off(document.body, 'click', fn) // no longer logs '!' upon clicking on the page

function off(el: HTMLElement, type: string, listener: () => void, capture = false) {
  el.removeEventListener(type, listener, capture)
}
