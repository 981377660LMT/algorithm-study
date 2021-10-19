function* indexOfSubstrings(str: string, searchValue: string) {
  let i = 0
  while (i < str.length) {
    const hit = str.indexOf(searchValue, i)
    if (hit !== -1) {
      yield hit
      i = hit + 1
    } else break
  }
}

console.log([...indexOfSubstrings('tiktok tok tok tik tok tik', 'tik')]) // [0, 15, 23]
;[...indexOfSubstrings('tutut tut tut', 'tut')] // [0, 2, 6, 10]
;[...indexOfSubstrings('hello', 'hi')] // []
