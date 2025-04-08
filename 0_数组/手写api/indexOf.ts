export {}

function myIndexOf(str: string, searchString: string, position?: number): number {
  if (position === undefined) position = 0
  position = Math.max(0, Math.min(position, str.length))
  if (searchString === '') return position
  if (searchString.length > str.length) return -1

  const table = buildKMPTable(searchString)
  let j = 0
  for (let i = position; i < str.length; i++) {
    while (j > 0 && str[i] !== searchString[j]) {
      j = table[j - 1]
    }
    j += +(str[i] === searchString[j])
    if (j === searchString.length) {
      return i - j + 1
    }
  }
  return -1
}

function buildKMPTable(searchString: string): number[] {
  const table = Array(searchString.length).fill(0)
  let len = 0
  for (let i = 1; i < searchString.length; i++) {
    while (len > 0 && searchString[i] !== searchString[len]) {
      len = table[len - 1]
    }
    len += +(searchString[i] === searchString[len])
    table[i] = len
  }
  return table
}

console.log(buildKMPTable('abcabc')) // [0, 0, 0, 1, 2, 3]
console.log(myIndexOf('abcabc', 'abc')) // 0
console.log(myIndexOf('abcabc', 'bca')) // -1
console.log(myIndexOf('abcabc', 'abc', 1)) // 3
console.log(myIndexOf('abcabc', 'abc', 0)) // 0
