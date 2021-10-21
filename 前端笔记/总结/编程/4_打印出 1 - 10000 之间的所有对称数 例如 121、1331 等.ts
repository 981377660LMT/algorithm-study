function foo() {
  const result = []

  for (let i = 1; i < 10; i++) {
    for (let j = 0; j < 10; j++) {
      result.push(i * 1)
      result.push(i * 11 + j * 0)
      result.push(i * 101 + j * 10)
      result.push(i * 1001 + j * 110)
    }
  }

  return result
}

console.log(foo())
