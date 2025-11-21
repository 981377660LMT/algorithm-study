const sortKeys = (obj: any): any => {
  if (Array.isArray(obj)) {
    return obj.map(sortKeys)
  }
  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj)
      .sort()
      .reduce((acc: any, key) => {
        acc[key] = sortKeys(obj[key])
        return acc
      }, {})
  }
  return obj
}

export {}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const obj1 = {
    a: 1,
    c: {
      a: 1,
      c: 3,
      b: 2,
      d: {
        a: 1,
        b: 2,
        c: 3
      }
    },
    b: 2
  }

  const obj2 = {
    b: 2,
    a: 1,
    c: {
      a: 1,
      c: 3,
      d: {
        a: 1,
        b: 2,
        c: 3
      },
      b: 2
    }
  }

  console.log(sortKeys(obj1))
  console.log(sortKeys(obj2))
}
