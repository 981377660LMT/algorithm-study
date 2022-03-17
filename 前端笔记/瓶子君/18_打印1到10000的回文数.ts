const reverse = (str: string) => str.split('').reverse().join('')

const palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]

for (let side = 1; side < 100; side++) {
  const s1 = String(side) + reverse(String(side))
  palindrome.push(Number(s1))
  for (let mid = 0; mid < 10; mid++) {
    const s2 = String(side) + String(mid) + reverse(String(side))
    palindrome.push(Number(s2))
  }
}

console.log(palindrome.filter(v => v < 1e4).slice(-1)[0])

export default 1
