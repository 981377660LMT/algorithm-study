# Monge矩阵
# 题意:给定一个H x W 的矩阵，
# 判断对于所有的1≤i1≤i2≤H,1≤j1≤j2≤W,
# 是否满足Ai1ji + Ai2ji2 ≤ Ai2j1 + Ai1j2
# 每当我们从Monge矩阵中挑出两行与两列，并且考虑行列交叉处的4个元素，
# 左上角与右下角的和小于等于左下角与右上角元素的和
# 1 ≤ H, W ≤ 1000

# !等价于A(i,j)+A(i+1,j+1)<=A(i+1,j)+A(i,j+1)


from typing import List


def solve(matrix: List[List[int]]) -> bool:
    """是否为Monge矩阵"""
    ROW, COL = len(matrix), len(matrix[0])
    for r in range(ROW - 1):
        for c in range(COL - 1):
            if matrix[r][c] + matrix[r + 1][c + 1] > matrix[r + 1][c] + matrix[r][c + 1]:
                return False
    return True


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]
    print("Yes" if solve(grid) else "No")
