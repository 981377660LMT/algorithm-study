onClickOutside(document.querySelector('#my-element')!, () => console.log('Hello'))

function onClickOutside(ele: Element, callback: () => void) {
  document.addEventListener('click', e => {
    if (!ele.contains(e.target as Node)) callback()
  })
}
// Will log 'Hello' whenever the user clicks outside of #my-element
