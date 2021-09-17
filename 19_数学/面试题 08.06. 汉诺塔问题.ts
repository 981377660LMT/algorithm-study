/**
 Do not return anything, modify C in-place instead.
 每一个盘子只能放在更大的盘子上面
 请编写程序，用栈将所有盘子从第一根柱子移到最后一根柱子。
 A:起点
 B；缓冲
 C：终点
 关注起点和终点即可
 */
function hanota(A: number[], B: number[], C: number[]): void {
  const move = (n: number, A: number[], B: number[], C: number[]) => {
    if (n === 1) return C.push(A.pop()!)
    move(n - 1, A, C, B)
    C.push(A.pop()!)
    move(n - 1, B, A, C)
  }
  move(A.length, A, B, C)
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
