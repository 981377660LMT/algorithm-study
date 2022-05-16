function numSteps(s: string) {
  let num = BigInt('0b' + s)
  let step = 0

  while (num !== 1n) {
    if (num % 2n === 0n) {
      num /= 2n
    } else {
      num += 1n
    }

    step++
  }

  return step
}

function numSteps2(s: string) {
  let res = 0
  let carry = 0

  // 注意i>0
  for (let i = s.length - 1; i > 0; i--) {
    if (s[i] === '0') {
      res += 1 + carry
    } else {
      res += 2 - carry
      carry = 1
    }
  }

  return res + carry
}

console.log(numSteps('1011111111010101010011'))
console.log(BigInt('0x1010'))
console.log(BigInt('0b1010'))
console.log(BigInt('0o1010'))
