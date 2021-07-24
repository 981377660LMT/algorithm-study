const isRotation = (s1: string, s2: string) => {
  if (s1.length !== s2.length) return false

  for (let index = 0; index < s1.length; index++) {
    if (s1 === s2) return true
    const s1Arr = s1.split('')
    s1Arr.push(s1Arr.shift()!)
    s1 = s1Arr.join('')
  }

  return false
}

console.log(isRotation('waterbottle', 'erbottlewat'))
console.log(isRotation('apple', 'ppale'))

export {}
