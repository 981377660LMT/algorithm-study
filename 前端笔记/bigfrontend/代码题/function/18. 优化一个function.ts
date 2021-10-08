/**
 * @param {object[]} items
 * @excludes { Array< {k: string, v: any} >} excludes
 */

/**
 * @param {object[]} items
 * @param { Array< {k: string, v: any} >} excludes
 * @return {object[]}
 */
function excludeItems(
  items: Record<string, any>[],
  excludes: { k: string; v: any }[]
): Record<string, any>[] {
  const excludesMap = new Map<string, Set<string>>()

  for (const { k, v } of excludes) {
    !excludesMap.has(k) && excludesMap.set(k, new Set())
    excludesMap.get(k)!.add(v)
  }

  return items.filter(
    item =>
      !Object.keys(item).some(key => excludesMap.has(key) && excludesMap.get(key)!.has(item[key]))
  )
}

let items = [
  { color: 'red', type: 'tv', age: 18 },
  { color: 'silver', type: 'phone', age: 20 },
  { color: 'blue', type: 'book', age: 17 },
]

// 一个由key和value组成的array
const excludes = [
  { k: 'color', v: 'silver' },
  { k: 'type', v: 'tv' },
]

function rawExcludeItems(items: object[], excludes: { k: string; v: any }[]) {
  excludes.forEach(pair => {
    // @ts-ignore
    items = items.filter(item => item[pair.k] === item[pair.v])
  })

  return items
}
