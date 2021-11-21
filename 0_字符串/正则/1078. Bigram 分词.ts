function findOcurrences(text: string, first: string, second: string): string[] {
  const reg = new RegExp(`(?<=\\b${first}\\s${second}\\s)\\w+`, 'g')
  return text.match(reg) ?? []
}

console.log(findOcurrences('alice is a good girl she is a good student', 'a', 'good'))

function findOcurrences2(text: string, first: string, second: string): string[] {
  const res: string[] = []
  const arr = text.split(' ')

  for (let i = 0; i < arr.length - 2; i++) {
    if (arr[i] === first && arr[i + 1] === second) {
      res.push(arr[i + 2])
    }
  }

  return res
}
