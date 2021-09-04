const plus = () => {
  let count = 0

  const plus = (...nums: number[]) => {
    nums.forEach(num => (count += num))
    return plus
  }

  plus.toString = () => count
  // plus.valueOf = () => count

  return plus
}

const add = plus()
console.log(`${add(1)(2, 3)}`)
export {}
