// 字符串仅由小写字母和[]组成

interface Node<V> {
  value: V
  children?: Node<V>
}

const normalize = (str: string): Node<string> => {
  const keys = str.match(/\w+/g)
  console.log(keys)
  if (!keys) throw new Error('invald input')
  if (keys.length === 1) return { value: keys[0] }

  const res = { root: {} } as Record<string, any>

  let root = res.root
  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i]
    root.value = key
    root.children = {}
    root = root.children
  }

  root.value = keys[keys.length - 1]

  return res.root as Node<string>
}

console.log(normalize('[abc[bcd[def]]]'))
console.log(normalize('abc'))

export {}
