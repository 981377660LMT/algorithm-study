/**
 *
 * @param start
 * @param end
 * @returns 返回 [start,end] 闭区间里的随机整数
 */
const randint = (start: number, end: number) => {
  if (start > end) throw new Error('invalid interval')
  const diff = end - start
  return Math.floor((diff + 1) * Math.random()) + start
}

if (require.main === module) {
  console.log(randint(1, 3))
}

export { randint }
