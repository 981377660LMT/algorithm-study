function shift(arr) {
  let len = arr.length
  if (len === 0) {
    return
  }
  let first = arr[0]
  for (let i = 0; i < len - 1; i++) {
    arr[i] = arr[i + 1]
  }
  arr.length = len - 1
  return first
}
