export {}
const x = [
  { id: 1, name: 'John' },
  { id: 2, name: 'Maria' },
]
const y = [{ id: 1, age: 28 }, { id: 3, age: 26 }, { age: 3 }]
console.log(combine(x, y, 'id'))

function combine(arr1: Record<any, any>[], arr2: Record<any, any>[], prop: string) {
  return Object.values(
    [...arr1, ...arr2].reduce((obj, cur) => {
      if (cur[prop]) {
        const key = cur[prop]
        obj[key] = obj[key] ? { ...obj[key], ...cur } : { ...cur }
      }
      return obj
    }, {})
  )
}
