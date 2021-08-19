/**
 * @param {number} n
 * @param {number} m
 * @return {number}
 * 每次从这个圆圈里删除第m个数字（删除后从下一个数字开始计数）。求出这个圆圈里剩下的最后一个数字。
 * 2.数组模拟
 * 数组模拟成循环链表
 * index=(index+m-1)%(list.size());
 * 但是还可能存在n本身就很大的情况，无论是顺序表ArrayList还是链表LinkedList去频繁查询、删除都是很低效的。
 */
const lastRemaining = (n: number, m: number): number => {
  if (m === 1) return n - 1
  const list: number[] = Array.from<number, number>({ length: n }, (_, i) => i)

  let index = 0
  while (list.length >= 2) {
    // 删除向后走的第m-1个元素
    index = (index + m - 1) % list.length
    list.splice(index, 1)
  }

  return list[0]
}

console.log(lastRemaining(5, 3))
// console.log(lastRemaining(10, 17))

export default 1
