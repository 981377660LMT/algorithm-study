// myJoin

function join<T>(arr: ArrayLike<T>, separator: T): T[] {
  const res: T[] = []
  for (let i = 0; i < arr.length; i++) {
    if (res.length) res.push(separator)
    res.push(arr[i])
  }
  return res
}

export { join }
