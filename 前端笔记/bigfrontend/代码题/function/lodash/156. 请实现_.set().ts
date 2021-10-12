// 由于不需要关心属性的存在与否，_.set(object, path, value) 非常方便实用。
// 请自行实现set()。

function set(obj: Record<PropertyKey, any>, path: string | string[], value: any): void {
  // your code here
  if (typeof path === 'string') path = normalize(path)
  if (path.length === 0) return

  let root = obj
  path.forEach((key, index, arr) => {
    if (index === path.length - 1) root[key] = value
    else {
      if (!(key in root)) {
        const next = arr[index + 1]
        root[key] = String(parseInt(next)) === next ? [] : {}
      }
      root = root[key]
    }
  })
}

function normalize(path: string): string[] {
  return path.split(/[\.\[\]]/g).filter(s => s.trim())
}

if (require.main === module) {
  const obj = {
    a: {
      b: {
        c: [1, 2, 3],
      },
    },
  }

  set(obj, 'a.b.c', 'BFE')
  console.log(obj.a.b.c) // "BFE"

  set(obj, 'a.b.c.0', 'BFE')
  console.log(obj.a.b.c[0]) // "BFE"

  set(obj, 'a.b.c[1]', 'BFE')
  console.log(obj.a.b.c[1]) // "BFE"

  set(obj, ['a', 'b', 'c', '2'], 'BFE')
  console.log(obj.a.b.c[2]) // "BFE"

  set(obj, 'a.b.c[3]', 'BFE')
  console.log(obj.a.b.c[3]) // "BFE"

  // set(obj, 'a.c.d[0]', 'BFE')
  // // valid digits treated as array elements
  // // @ts-ignore
  // console.log(obj.a.c.d[0]) // "BFE"

  // set(obj, 'a.c.d.01', 'BFE')
  // // invalid digits treated as property string
  // // @ts-ignore
  // console.log(obj.a.c.d['01']) // "BFE"
}

export {}
