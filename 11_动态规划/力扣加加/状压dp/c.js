const a = [12, 3]
const b = [5, 77]
console.log(
  a.reduce((pre, cur, index, arr) => {
    console.log(index)
    return pre + arr[index] * b[index]
  }, 0)
)
