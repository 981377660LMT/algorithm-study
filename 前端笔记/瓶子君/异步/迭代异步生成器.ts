async function* asyncGenerator() {
  let i = 0
  while (i < 3) {
    yield i++
  }
}

;(async function () {
  for await (const g of asyncGenerator()) {
    console.log(g)
  }
})()

export {}
