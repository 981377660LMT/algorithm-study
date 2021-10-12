// fetchList is provided for you
declare const fetchList: (since?: number) => Promise<{ items: { id: number }[] }>

// you can change this to generator function if you want
const fetchListWithAmount = async (amount = 5) => {
  // your code here
  const result: { id: number }[] = []

  while (result.length < amount) {
    const lastItem = result[result.length - 1]
    const { items } = await fetchList(lastItem?.id)
    if (items.length > 0) {
      result.push(...items)
    } else break // 服务器已经没有数据了
  }

  return result.slice(0, amount)
}

if (require.main === module) {
  fetchListWithAmount(10).then(dataGet => {
    expect(dataGet).toEqual([1, 2, 3, 4, 5, 6, 7, 8, 10, 11].map(item => ({ id: item })))
  })
}
