// // _.get(object, path, [defaultValue]) 是一个很好用的方法，用来从复杂数据结构中获取特定属性的值。
// // 第三个参数是目标值为undefined的时候的默认值。

/**
 * @param {object} source
 * @param {string | string[]} path
 * @param {any} [defaultValue]
 * @return {any}
 */
function get(
  source: Record<PropertyKey, any>,
  path: string | string[],
  defaultValue: any = undefined
): any {
  if (typeof path === 'string') path = normalize(path)
  if (path.length === 0) return defaultValue
  let root = source
  for (const key of path) {
    if (!(key in root)) return defaultValue
    root = root[key]
  }
  return root
}

if (require.main === module) {
  const obj = {
    a: {
      b: {
        c: [1, 2, 3],
      },
    },
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

function normalize(path: string): string[] {
  return path.split(/[\.\[\]]/g).filter(s => s.trim())
}

// console.log(normalize('a.b.c[1]'))
