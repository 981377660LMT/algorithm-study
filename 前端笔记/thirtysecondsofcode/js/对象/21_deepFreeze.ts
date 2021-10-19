const myObj = {
  a: 1,
  b: 'hello',
  c: [0, 1, 2],
  d: { e: 1, f: 2 },
}

const deepFreeze = (obj: Record<any, any>) => {
  for (const prop of Object.keys(obj)) {
    if (obj[prop] == null) continue
    if (typeof obj[prop] === 'object' && !Object.isFrozen(obj[prop])) deepFreeze(obj[prop])
  }

  return Object.freeze(obj)
}

deepFreeze(myObj)

myObj.a = 10
myObj.b = 'hi'
myObj.c[1] = 4
myObj.d.e = 0

// console.log(myObj)

// console.log(Object.keys(null))  报错
