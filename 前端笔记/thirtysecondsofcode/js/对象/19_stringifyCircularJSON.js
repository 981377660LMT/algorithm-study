const obj = { n: 42 }
obj.obj = obj
console.log(stringifyCircularJSON(obj)) // '{"n": 42}'

function stringifyCircularJSON(obj) {
  const visited = new WeakSet()

  return JSON.stringify(obj, (key, value) => {
    if (value != null && typeof value == 'object') {
      if (visited.has(value)) return
      visited.add(value)
    }

    return value
  })
}
