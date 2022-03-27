let line

while ((line = readline()) != '') {
  const arr = line.split(' ')
  let n = Number(arr[0])
  const m = Number(arr[1])
  let res = 0

  for (let _ = 0; _ < m; _++) {
    res += n
    n = Math.sqrt(n)
  }

  print(res.toFixed(2))
}
