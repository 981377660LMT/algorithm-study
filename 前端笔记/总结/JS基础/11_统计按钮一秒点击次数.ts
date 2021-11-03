const button = document.getElementById('button')!
button.onclick = countClickInInterval()

function countClickInInterval(ms = 1000) {
  let count = 0
  let timer: number | undefined

  function inner() {
    count++
    if (timer != undefined) return

    timer = window.setTimeout(() => {
      console.log(`点击了${count}次`)
      count = 0
      timer = undefined
    }, ms)
  }

  return inner
}
