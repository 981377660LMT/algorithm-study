# 给定一个n * m 的矩阵w,行和列都从1开始,
# 对于(i,j),如果i+j是偶数那么i,j=(i-1) * m +j,否则wi,j =0。
# 现有q个查询,每次询问一个子矩阵的和,其中左上角为(x1,y1),右下角矩阵为(x2, y2)。
# 数据范围:
# 1 ≤n, m ≤1e9,q≤ 2* 1e5
# 1 ≤x1, x2≤n,1≤ y1, y2 ≤m


# 矩阵为
# 1 2 3 4 5

# 如果是(r+c)是奇数 则替换为0
# 求子矩阵的和
# !推公式要静下心来推 多写一些辅助函数简化逻辑
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def arithmeticSum1(first: int, last: int, diff: int) -> int:
        """等差数列求和 first:首项 last:末项 diff:公差"""
        item = (last - first) // diff + 1
        return item * (first + last) // 2

    def arithmeticSum2(first: int, diff: int, item: int) -> int:
        """等差数列求和 first:首项 diff:公差 item:项数"""
        last = first + (item - 1) * diff
        return item * (first + last) // 2

    def cal(row: int, col: int) -> int:
        """(0,0) 到 (row,col)的子矩阵的和

        分成奇数行/偶数行的等差数列和
        """
        if row == 0 or col == 0:
            return 0
        row, col = row - 1, col - 1

        res = 0

        # 偶数行
        first, diff1, item1 = 1, 2, col // 2 + 1
        firstRowSum = arithmeticSum2(first, diff1, item1)
        diff2, item2 = 2 * COL * item1, row // 2 + 1
        allSum = arithmeticSum2(firstRowSum, diff2, item2)
        res += allSum

        # 奇数行
        if col >= 1 and row >= 1:
            first, diff1, item1 = COL + 2, 2, (col + 1) // 2
            firstRowSum = arithmeticSum2(first, diff1, item1)
            diff2, item2 = 2 * COL * item1, (row + 1) // 2
            allSum = arithmeticSum2(firstRowSum, diff2, item2)
            res += allSum

        return res % MOD

    ROW, COL = map(int, input().split())
    q = int(input())
    for _ in range(q):
        row1, row2, col1, col2 = map(int, input().split())
        row1, row2, col1, col2 = row1 - 1, row2 - 1, col1 - 1, col2 - 1
        print(
            (cal(row2 + 1, col2 + 1) - cal(row1, col2 + 1) - cal(row2 + 1, col1) + cal(row1, col1))
            % MOD
        )
