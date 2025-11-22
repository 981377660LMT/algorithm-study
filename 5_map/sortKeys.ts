const deepSortKeys = (json: any): any => {
  if (Array.isArray(json)) {
    return json.map(deepSortKeys)
  }
  if (json !== null && typeof json === 'object') {
    return Object.keys(json)
      .sort()
      .reduce((acc: any, key) => {
        acc[key] = deepSortKeys(json[key])
        return acc
      }, {})
  }
  return json
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
        d: [
          2,
          {
            a: 1
          },
          3
        ],
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

  console.dir(deepSortKeys(obj1), { depth: null })
  console.dir(deepSortKeys(obj2), { depth: null })
}
