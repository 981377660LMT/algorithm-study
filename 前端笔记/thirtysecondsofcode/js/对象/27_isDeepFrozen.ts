const x = Object.freeze({ a: 1 })
const y = Object.freeze({ b: { c: 2 } })
isDeepFrozen(x) // true
isDeepFrozen(y) // false

function isDeepFrozen(obj: Record<any, any>): boolean {
  return (
    Object.isFrozen(obj) &&
    Object.keys(obj).every(prop => typeof obj[prop] !== 'object' || isDeepFrozen(obj[prop]))
  )
}

export {}
