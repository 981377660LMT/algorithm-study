var convertToBaseN = function (num: number, n: number): string {
  if (num < 0) return `-${convertToBaseN(num * -1, n)}`
  if (num < n) return `${n}`
  return convertToBaseN(num / n, n) + convertToBaseN(num % n, n)
}