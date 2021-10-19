const orderWith = (arr: Record<string, any>[], prop: string, order: any[]) => {
  const orderValues = order.reduce((pre, cur, index) => {
    pre[cur] = index
    return pre
  }, {})

  return [...arr].sort((a, b) => {
    if (orderValues[a[prop]] === undefined) return 1
    if (orderValues[b[prop]] === undefined) return -1
    return orderValues[a[prop]] - orderValues[b[prop]]
  })
}
const users = [
  { name: 'fred', language: 'Javascript' },
  { name: 'barney', language: 'TypeScript' },
  { name: 'frannie', language: 'Javascript' },
  { name: 'anna', language: 'Java' },
  { name: 'jimmy' },
  { name: 'nicky', language: 'Python' },
]
orderWith(users, 'language', ['Javascript', 'TypeScript', 'Java'])
