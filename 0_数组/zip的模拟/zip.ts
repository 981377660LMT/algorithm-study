/**
 *
 * @param strs
 * @returns 将每个列组合 则有最大长度的行，原来的单词数的列
 */
// function zipLongest(fillValue = '', ...strs: string[]) {
//   const length = Math.max(...strs.map(str => str.length))
//   // 行对列，列对行
//   const arr = Array.from({ length }, () => Array(strs.length).fill(''))
//   for (let i = 0; i < strs.length; i++) {
//     for (let j = 0; j < length; j++) {
//       arr[j][i] = strs[i][j] || fillValue
//     }
//   }
//   return arr
// }

function* zip<T>(...arr: ArrayLike<T>[]) {
  const length = Math.min(...arr.map(arrlike => arrlike.length))
  for (let i = 0; i < length; i++) {
    yield arr.map(arrlike => arrlike[i])
  }
}

function* zipLongest(...arr: ArrayLike<any>[]) {
  const length = Math.max(...arr.map(arrlike => arrlike.length))
  for (let i = 0; i < length; i++) {
    yield arr.map(arrlike => arrlike[i])
  }
}

export { zip, zipLongest }

if (require.main === module) {
  console.log(...zipLongest('12', [1, 2, 3], [90, 8]))
}
