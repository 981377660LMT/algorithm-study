const result: number[] = []

// A 型 AA 型 ABA型 ABBA型
for (let i = 1; i < 10; i++) {
  result.push(i)
  result.push(i * 11)
  for (let j = 0; j < 10; j++) {
    result.push(i * 101 + j * 10)
    result.push(i * 1001 + j * 110)
  }
}

console.log(result)

export default 1
