import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 矩阵为
# 1 2 3 4 5

# 如果是(r+c)是奇数 则替换为0
# 求子矩阵的和
# 数学太差不会推
if __name__ == "__main__":

    def cal(row: int, col: int) -> int:
        """(0,0) 到 (row,col)的子矩阵的和"""
        if row == 0 or col == 0:
            return 0
        row, col = row - 1, col - 1

    n, m = map(int, input().split())
    q = int(input())
    for _ in range(q):
        row1, row2, col1, col2 = map(int, input().split())
        row1, row2, col1, col2 = row1 - 1, row2 - 1, col1 - 1, col2 - 1
        print(
            (cal(row2 + 1, col2 + 1) - cal(row1, col2 + 1) - cal(row2 + 1, col1) + cal(row1, col1))
            % MOD
        )
