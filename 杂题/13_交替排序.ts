const wiggle = (arr: number[]) => {
  arr.sort((a, b) => a - b)
  const mid = Math.floor(arr.length / 2)
  const odd = arr.slice(0, mid)
  const even = arr.slice(mid)

  return even.reduce<number[]>((pre, cur, index) => {
    pre.push(cur)
    odd[index] && pre.push(odd[index])
    return pre
  }, [])
}

console.log(wiggle([2, 1, 3, 4, 5, 6, 7]))

export {}
