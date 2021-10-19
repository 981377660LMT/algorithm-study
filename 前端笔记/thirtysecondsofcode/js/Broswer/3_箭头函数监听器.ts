const toggleElements = document.querySelectorAll('.toggle')
toggleElements.forEach(el => {
  el.addEventListener('click', () => {
    this.classList.toggle('active') // `this` refers to `window`
    // Error: Cannot read property 'toggle' of undefined
  })
})

const toggleElements2 = document.querySelectorAll('.toggle')
toggleElements2.forEach(el => {
  ;(el as HTMLButtonElement).addEventListener('click', e => {
    e.currentTarget!.classList.toggle('active') // works correctly
  })
})
