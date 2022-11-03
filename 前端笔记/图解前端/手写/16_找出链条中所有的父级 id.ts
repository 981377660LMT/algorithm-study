interface ListItem {
  id: string
  name: string
  children?: ListItem[]
}

const list: ListItem[] = [
  {
    id: '1',
    name: 'test1',
    children: [
      {
        id: '11',
        name: 'test11',
        children: [
          {
            id: '111',
            name: 'test111'
          },
          {
            id: '112',
            name: 'test112'
          }
        ]
      },
      {
        id: '12',
        name: 'test12',
        children: [
          {
            id: '121',
            name: 'test121'
          },
          {
            id: '122',
            name: 'test122'
          }
        ]
      }
    ]
  }
]
const id = '112'

console.log(fn(id, list)) // 输出 [1， 11， 112]

function fn(id: string, list: ListItem[]) {
  const path: string[] = []

  const bt = (curItem: ListItem[], index: number, path: string[]) => {
    console.log(index)
    if (index === id.length) return path
    const next = curItem.find(item => id.startsWith(item.id))
    if (!next) return
    path.push(next.id)
    bt(next.children!, index + 1, path)
  }

  bt(list, 0, path)
  return path
}

export {}
