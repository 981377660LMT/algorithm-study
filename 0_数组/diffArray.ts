// 求两个数组的 changetset.
// 使用 lodash 以及不使用 lodash 两种写法.

import { difference, intersection, isEqual } from 'lodash-es'

function diffArray1<T>(
  pre: T[],
  cur: T[],
  options?: {
    key?: (item: T) => PropertyKey
    equals?: (a: T, b: T) => boolean
  }
) {
  const key = options?.key ?? ((item: any) => item)
  const equals = options?.equals ?? ((a, b) => a === b)

  const preMap = new Map<PropertyKey, T>()
  pre.forEach(item => {
    preMap.set(key(item), item)
  })

  const curMap = new Map<PropertyKey, T>()
  cur.forEach(item => {
    curMap.set(key(item), item)
  })

  const added: T[] = []
  const removed: T[] = []
  const updated: T[] = []

  preMap.forEach((value, key) => {
    if (!curMap.has(key)) {
      removed.push(value)
    }
  })

  curMap.forEach((value, key) => {
    if (!preMap.has(key)) {
      added.push(value)
    } else {
      const preValue = preMap.get(key)!
      if (!equals(preValue, value)) {
        updated.push(value)
      }
    }
  })

  return {
    added,
    removed,
    updated
  }
}

function diffArray2<T>(
  pre: T[],
  cur: T[],
  options?: {
    key?: (item: T) => PropertyKey
    equals?: (a: T, b: T) => boolean
  }
) {
  const key = options?.key ?? (item => item)
  const equals = options?.equals ?? ((a, b) => a === b)

  const preMap = new Map(pre.map(item => [key(item), item]))
  const curMap = new Map(cur.map(item => [key(item), item]))
  const preKeys = [...preMap.keys()]
  const curKeys = [...curMap.keys()]

  const removedKeys = difference(preKeys, curKeys)
  const addedKeys = difference(curKeys, preKeys)
  const commonKeys = intersection(preKeys, curKeys)

  const removed = removedKeys.map(key => preMap.get(key)!)
  const added = addedKeys.map(key => curMap.get(key)!)
  const updated = commonKeys
    .filter(key => !equals(preMap.get(key)!, curMap.get(key)!))
    .map(key => curMap.get(key)!)

  return {
    added,
    removed,
    updated
  }
}

export {}

console.log(
  diffArray1(
    [
      { id: 1, name: 'a' },
      { id: 2, name: 'b' }
    ],
    [
      { id: 1, name: 'a2' },
      { id: 3, name: 'c' }
    ],
    {
      key: item => item.id,
      equals: (a, b) => isEqual(a, b)
    }
  )
)
