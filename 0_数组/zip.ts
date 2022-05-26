function* zip<T>(...arr: ArrayLike<T>[]) {
  const length = Math.min(...arr.map(arrlike => arrlike.length))
  for (let i = 0; i < length; i++) {
    yield arr.map(arrlike => arrlike[i])
  }
}

function* zipLongest<T>(
  fillValue: any = undefined,
  ...arr: ArrayLike<T>[]
): Generator<T[], void, unknown> {
  const length = Math.max(...arr.map(arrlike => arrlike.length))
  for (let i = 0; i < length; i++) {
    yield arr.map(arrlike => arrlike[i] ?? fillValue)
  }
}

export { zip, zipLongest }

if (require.main === module) {
  console.log(...zipLongest('12', [1, 2, 3], [90, 8]))
}
