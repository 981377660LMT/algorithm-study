var inputArray = [1, 2, 3]

var result = inputArray.reduce(function permute(pre, cur, index, arr) {
  return pre.concat(
    (arr.length > 1 &&
      arr
        .slice(0, index)
        .concat(arr.slice(index + 1))
        .reduce(permute, [])
        .map(perm => [cur].concat(perm))) ||
      cur
  )
}, [])

console.log(result)
