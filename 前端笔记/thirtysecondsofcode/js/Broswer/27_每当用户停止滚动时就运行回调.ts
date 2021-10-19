onScrollStop(() => {
  console.log('The user has stopped scrolling')
})

function onScrollStop(callback: () => void) {
  let isScrolling: NodeJS.Timer
  window.addEventListener(
    'scroll',
    e => {
      clearTimeout(isScrolling)
      isScrolling = setTimeout(() => {
        callback()
      }, 150)
    },
    false
  )
}

export {}
