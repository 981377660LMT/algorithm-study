// 实现loadsh的_.get()方法，用来从复杂数据结构中获取特定属性的值。
// _.get(object, path, [defaultValue]) 是一个很好用的方法，用来从复杂数据结构中获取特定属性的值。
// 第三个参数是目标值为undefined的时候的默认值。

function get(
  source: Record<PropertyKey, unknown>,
  path: string | string[],
  defaultValue?: unknown
): unknown {
  if (typeof path === 'string') path = normalize(path)
  if (path.length === 0) return defaultValue

  let root: unknown = source
  for (const key of path) {
    if (!isRecord(root)) return defaultValue
    if (!(key in root)) return defaultValue // Object.hasOwnProperty()
    root = root[key]
  }

  return root
}

function isRecord(value: unknown): value is Record<PropertyKey, unknown> {
  return typeof value === 'object' && value !== null
}

function normalize(path: string): string[] {
  return path.split(/\.|\[|\]/g).filter(s => s.trim())
}

if (require.main === module) {
  const obj = {
    a: {
      b: {
        c: [1, 2, 3]
      }
    }
  }

  console.log(get(obj, 'a.b.c')) // [1,2,3]
  console.log(get(obj, 'a.b.c.0')) // 1
  console.log(get(obj, 'a.b.c[1]')) // 2
  console.log(get(obj, ['a', 'b', 'c', '2'])) // 3
  console.log(get(obj, 'a.b.c[3]')) // undefined
  console.log(get(obj, 'a.c', 'bfe')) // 'bfe'
  console.log(get(obj, [])) // undefined
  console.log(get(obj, '')) // undefined
}

export {}
