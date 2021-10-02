/**
 *
 * @param num  num * num 的二维矩阵
 * @param xPos
 * @param yPos
 * 乐团共有 9 种乐器，乐器编号为 1~9，每位成员持有 1 个乐器。
 * https://leetcode-cn.com/problems/SNJvJP/solution/gen-ju-mian-ji-qiu-chai-zhao-dao-suo-zai-06ih/
 */
function orchestraLayout(num: number, xPos: number, yPos: number): number {
  // 处于第几层 从0开始算
  const layer = Math.min(xPos, num - 1 - xPos, yPos, num - 1 - yPos)
  // 前几圈有多少个元素  模9即当前圈左上角编号
  let index = (num ** 2 - (num - 2 * layer) ** 2) % 9
  // (start, end 从左上角开始计算)
  const [start, end] = [layer, num - layer - 1]
  const edge = end - start

  // 上
  if (xPos === start) return ((index + yPos - start) % 9) + 1
  // 右
  else if (yPos === end) return ((index + edge + xPos - start) % 9) + 1
  // 下
  else if (xPos === end) return ((index + edge * 2 + end - yPos) % 9) + 1
  // 左
  else if (yPos === start) return ((index + edge * 3 + end - xPos) % 9) + 1

  throw Error('invalid input')
}

console.log(orchestraLayout(4, 1, 2))
console.log(orchestraLayout(3, 0, 2))
console.log(orchestraLayout(2511, 1504, 2235)) // 3
// 自 grid 左上角开始顺时针螺旋形向内循环以 1，2，...，9 循环重复排列
// 请返回位于场地坐标 [Xpos,Ypos] 的成员所持乐器编号

// 力扣标记的简单题不一定简单，但是困难是真的困难

// 怕了
