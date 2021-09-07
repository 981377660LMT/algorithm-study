/**
 * @param {string[]} list1
 * @param {string[]} list2
 * @return {string[]}
 * 你需要帮助他们用最少的索引和找出他们共同喜爱的餐厅
 * @summary 类似于两数之和的做法
 */
const findRestaurant = function (list1: string[], list2: string[]): string[] {
  const res: [number, string][] = []
  const map = new Map<string, number>()
  for (let i = 0; i < list1.length; i++) {
    map.set(list1[i], i)
  }

  let min = Infinity
  for (let i = 0; i < list2.length; i++) {
    if (!map.has(list2[i])) continue
    const indexSum = i + map.get(list2[i])!
    const place = list2[i]
    res.push([indexSum, place])
    min = Math.min(min, indexSum)
  }

  return res.filter(v => v[0] === min).map(v => v[1])
}

console.log(
  findRestaurant(
    ['Shogun', 'Tapioca Express', 'Burger King', 'KFC'],
    ['Piatti', 'The Grill at Torrey Pines', 'Hungry Hunter Steakhouse', 'Shogun']
  )
)

export {}
