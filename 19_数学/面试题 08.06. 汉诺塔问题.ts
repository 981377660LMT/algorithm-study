/**
 Do not return anything, modify C in-place instead.
 每一个盘子只能放在更大的盘子上面
 请编写程序，用栈将所有盘子从第一根柱子移到最后一根柱子。
 A:起点
 B；缓冲
 C：终点
 关注起点和终点即可
 */
function hanota(from: number[], tmp: number[], to: number[]): void {
  const move = (n: number, from: number[], tmp: number[], to: number[]) => {
    if (n === 1) {
      return to.push(from.pop()!)
    }
    move(n - 1, from, to, tmp) // 第一个移到中间
    to.push(from.pop()!) // 最大的移到第三个
    move(n - 1, tmp, from, to) // 中间的移到第三个
  }
  move(from.length, from, tmp, to)
}
const C: number[] = []
console.log(hanota([2, 1, 0], [], C))
// 假设有n个盘子需要移动
// 首先将最上面的n-1个盘子从A移到B柱子
// 然后将最下面的一个盘子从A移到C柱子
// 最后将n-1个盘子从B移到C柱子
// 以上是汉诺塔的整体操作，其中移动n-1个盘子的操作是递归操作
console.log(C)
export {}
function hanota2(from: string, tmp: string, to: string): void {
  const move = (n: number, from: string, tmp: string, to: string) => {
    if (n === 1) {
      return console.log(`from ${from} to ${to}`)
    }
    move(n - 1, from, to, tmp) // 第一个移到中间
    move(1, from, tmp, to)
    move(n - 1, tmp, from, to) // 中间的移到第三个
  }
  move(3, from, tmp, to)
}

hanota2('a', 'b', 'c')
