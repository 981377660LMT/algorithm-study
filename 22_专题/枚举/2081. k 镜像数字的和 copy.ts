const reverse = (str: string) => str.split('').reverse().join('')

const palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]

for (let side = 1; side < 100000; side++) {
  const s1 = String(side) + reverse(String(side))
  palindrome.push(Number(s1))
  for (let mid = 0; mid < 10; mid++) {
    const s2 = String(side) + String(mid) + reverse(String(side))
    palindrome.push(Number(s2))
  }
}

palindrome.sort((a, b) => a - b)

function kMirror(k: number, n: number): number {
  const res: number[] = []
  let index = 0

  while (res.length < n) {
    const cand = palindrome[index].toString(k)
    if (cand === reverse(cand)) res.push(palindrome[index])
    index++
  }

  return res.reduce((pre, cur) => pre + cur, 0)
}

console.log(kMirror(7, 17))
console.log(kMirror(3, 12)) // 31730

export {}
