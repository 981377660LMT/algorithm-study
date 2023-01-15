// 实现lodash的set()方法，该方法可以对嵌套的对象进行赋值。
// 由于不需要关心属性的存在与否，_.set(object, path, value) 非常方便实用。
// 请自行实现set()。

function set(obj: Record<PropertyKey, unknown>, path: string | string[], value: unknown): void {
  path = normalize(path)
  let root = obj
  for (let i = 0; i < path.length; i++) {
    const key: string = path[i]
    if (i === path.length - 1) {
      root[key] = value
      continue
    }

    if (!Object.prototype.hasOwnProperty.call(root, key)) {
      const nextKey = path[i + 1]
      root[key] = isValidNumeric(nextKey) ? [] : {}
    }

    root = root[key] as Record<PropertyKey, unknown>
  }
}

function normalize(path: string | string[]): string[] {
  if (Array.isArray(path)) return path
  return path.split(/\.|\[|\]/g).filter(s => s.trim())
}

function isValidNumeric(str: string): boolean {
  if (str.length > 1 && str[0] === '0') return false
  return isDigit(str)
}

function isDigit(str: string): boolean {
  for (let i = 0; i < str.length; i++) {
    const char = str[i]
    if (
      !(
        char === '0' ||
        char === '1' ||
        char === '2' ||
        char === '3' ||
        char === '4' ||
        char === '5' ||
        char === '6' ||
        char === '7' ||
        char === '8' ||
        char === '9'
      )
    ) {
      return false
    }
  }

  return true
}

if (require.main === module) {
  const obj = {
    a: {
      b: {
        c: [1, 2, 3]
      }
    }
  }

  // set(obj, 'a.b.c', 'BFE')
  // console.log(obj.a.b.c) // "BFE"

  // set(obj, 'a.b.c.0', 'BFE')
  // console.log(obj.a.b.c[0]) // "BFE"

  // set(obj, 'a.b.c[1]', 'BFE')
  // console.log(obj.a.b.c[1]) // "BFE"

  // set(obj, ['a', 'b', 'c', '2'], 'BFE')
  // console.log(obj.a.b.c[2]) // "BFE"

  // set(obj, 'a.b.c[3]', 'BFE')
  // console.log(obj.a.b.c[3]) // "BFE"
}

export {}
