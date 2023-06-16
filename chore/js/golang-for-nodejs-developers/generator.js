function* gen() {
  yield 1
  yield 2
}

const g = gen()

while (true) {
  const { value, done } = g.next()
  if (done) {
    break
  }
  console.log(value)
}

for (const v of gen()) {
  console.log(v)
}
