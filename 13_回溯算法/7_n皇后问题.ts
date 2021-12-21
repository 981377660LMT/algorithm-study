/**
 * @param {number} n
 * @return {string[][]}
 * @description 皇后彼此不能相互攻击，也就是说：任何两个皇后都不能处于同一条横行、纵行或斜线上。
 * @description 如何快速剪枝:// 将index行的皇后摆放在第i列，在每一行的0到n-1列上尝试 判断列和两条对角线是否存在元素
 */
const solveNQueens = function (n: number): string[][] {
  const res: string[][] = []

  // path表示每行把皇后存第几列，记录列数
  const bt = (path: number[], curRow: number) => {
    if (curRow === n) {
      res.push(path.map(col => '.'.repeat(col) + 'Q' + '.'.repeat(n - col - 1)))
      return
    }

    for (let i = 0; i < n; i++) {
      if (
        path.some((col, row) => col === i || col + row === i + curRow || col - row === i - curRow)
      )
        continue
      path.push(i)
      bt(path, curRow + 1)
      path.pop()
    }
  }
  bt([], 0)

  return res
}

console.dir(solveNQueens(9), { depth: null })
// 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
export {}

// 返回 n 皇后问题 不同的解决方案的数量
// 二进制为 1 代表不可放置
// x & -x ：得到最低位的 1  代表除最后一位 1 保留，其他位全部为 0
// x & (x-1)：清零最低位的 1  代表将最后一位 1 变成 0
// x & ((1 << n) - 1)：将 x 最高位至第 n 位(含)清零 即与上n个1
function totalNQueens(n: number): number {
  let res = 0
  const bt = (n: number, row: number, col: number, mainDiagonal: number, subDiagonal: number) => {
    if (row == n) return res++
    let validPosition = ~(col | mainDiagonal | subDiagonal) & ((1 << n) - 1)
    while (validPosition) {
      const nextPosition = validPosition & -validPosition // 选取右边最后一个1
      validPosition &= validPosition - 1 // 移除右边最后一个1
      bt(
        n,
        row + 1,
        col | nextPosition,
        (mainDiagonal | nextPosition) >> 1,
        (subDiagonal | nextPosition) << 1
      )
    }
  }
  bt(n, 0, 0, 0, 0)
  return res
}

// void dfs(int n, int row, int col, int ld, int rd) {
//   if (row >= n) { res++; return; }
//   int bits = ~(col | ld | rd) & ((1 << n) - 1);   // 1
//   while (bits > 0) {   // 2
//       int pick = bits & -bits; // 3
//       dfs(n, row + 1, col | pick, (ld | pick) << 1, (rd | pick) >> 1); //4
//       bits &= bits - 1; // 5
//   }
// }
// 选取核心代码，按照上面注释的数字依次说：

// 1.(1 << n) - 1 这个语句实际上生成了n个1.这里的1表示可以放置皇后
// （其实就是初始化了n个1，在不考虑皇后之间可以相互攻击的情况下，n个位置都可以放皇后）；~(col | ld | rd)这里的三个变量分别代表了列以及两个斜线的放置情况。
// 这里的1表示的是不能放置皇后(就是相应的列或斜线上已经放置过皇后了)，
// 这与之前 (1 << n) - 1生成的n个1是不同含义的。
// 因此bits = ~(col | ld | rd) & ((1 << n) - 1)表示的是考虑了相应列、
// 斜线后能放置皇后的位置。
// 举个例子：n=4时，初始化为1111，
// 表示此时4个位置都可以放皇后，但是和~(col | ld | rd)按位与后变为了0110，
// 表示此时只有第2个和第3个位置是可以放皇后的。
// 2.当bits>0时，说明bits中还有1存在，就说明遍历还没有完成。
// 而在之后的循环体中，每遍历bits中的一个1，就会将其清0，
// 这就是代码中注释部分5的语句。
// 3.这里的pick就是取出了最后一位1，
// 表示此时遍历的是这种情况。假设bits为0110，取出最后一位1后，就变为0010，
// 就是将皇后放在第3个位置。
// 4.这里是核心：row+1不难理解，
// 就是因为之前已经在row行放置了皇后了，现在应该搜索下一行可能的位置了。
// col | pick就是把目前所有放置皇后的列都计算出来了，
// 比如最开始计算时col是0000，pick是0010,那么col | pick就是0010，
// 意思就是第三列被放置过了。接着说，假设ld是0000，ld | pick就是0010，
// 左移1位后变成了0100，意思就是下一行的第二列也不要放皇后了，
// 因为在这一行的第三列我已经放过了，他们是位于一个斜线上的。
// (rd | pick) >> 1跟(ld | pick) << 1是一个含义，就不赘述了。
