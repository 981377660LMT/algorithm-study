const plus = () => {
  let count = 0

  const inner = (...nums: number[]) => {
    nums.forEach(num => (count += num))
    return inner
  }

  inner.toString = () => count
  // plus.valueOf = () => count

  return inner
}

const add = plus()
console.log(`${add(1)(2, 3)}`)
export {}
