// fetchList is provided for you
declare const fetchList: (since?: number) => Promise<{ items: { id: number }[] }>

/**
 *
 * @param amount
 * @returns
 * 第一个request，直接调用fetchList，从response中取得最后一个item的id - lastItemId
 * 调用fetchList(lastItemId)获取下一个response
 * API 一次只返回5个item，加上一些过滤，实际的返回值中可能比5更少。如果一个都没有返回的话，就意味着服务器已经没有可以返回的了，我们需要停止调用。
 * 请实现一个函数用来获取任意数量的item
 */
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
