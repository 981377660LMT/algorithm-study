function solve(s) {
  let d1 = 0
  let d2 = 0
  let d3 = 0
  let res = 0
  for (let i = 0; i < s.length; i++) {
    const char = s[i]

    if (char === '(') d1++
    else if (char === '[') d2++
    else if (char === '{') d3++
    else if (char === ')') d1--
    else if (char === ']') d2--
    else if (char === '}') d3--

    if (d1 === -1) {
      res++
      d1++
    }

    if (d2 === -1) {
      res++
      d2++
    }

    if (d3 === -1) {
      res++
      d3++
    }
  }

  return res + d1 + d2 + d3
}

console.log(solve('))]'))
console.log(solve('(([]](])'))
console.log(solve('(())'))
console.log(solve('([()])[]'))
console.log(solve('[))]'))
console.log(solve('[))'))
console.log(solve('[)}'))
export {}
