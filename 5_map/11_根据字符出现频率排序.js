/**
 * @param {string} s
 * @return {string}
 */
var frequencySort = function (s) {
  const map = new Map()
  for (const letter of s) {
    if (map.has(letter)) {
      const [l, count] = map.get(letter)
      map.set(letter, [l, count + 1])
    } else {
      map.set(letter, [letter, 1])
    }
  }

  return [...map.values()]
    .sort((a, b) => b[1] - a[1])
    .map(item => item[0].repeat(item[1]))
    .join('')
}

console.log(frequencySort('tree'))
