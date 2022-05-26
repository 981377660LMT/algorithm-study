const swap = (i: number, j: number) => {
  i ^= j
  j ^= i
  i ^= j
  return [i, j]
}
console.log(swap(1, 2))

export default 1
