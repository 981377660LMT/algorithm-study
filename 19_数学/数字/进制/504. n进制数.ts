/**
 * @param {number} num
 * @return {string}
 * 不借助api num转n(n<10)进制
 */
const convertToBaseN = function (num: number, n: number): string {
  if (num < 0) return `-${convertToBaseN(num * -1, n)}`
  if (num < n) return `${num}`
  return convertToBaseN(~~(num / n), n) + convertToBaseN(num % n, n)
}

console.log(convertToBaseN(-7, 7))

export {}
