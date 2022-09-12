# 希尔伯特曲线
# https://leetcode.cn/problems/Bec9GP/?envType=study-plan&id=didizhenti
# 左下角的方格编号为 1，左上角为 2，右上角为 3，右下角为 4，
# 其余的图像也可以以此找到遍历网格的顺序。
# 下面我们希望给定任何边长为 2 的乘方的网格矩形，
# 给出基于 Hilbert 曲线索引的网格顺序，全部以左下角为第一个被索引的元素，
# 其索引值为 1，输出完整矩阵的索引值
# N <= 256

from typing import List


def solve(n) -> List[List[int]]:
    if n == 1:
        return [[1]]
    if n == 2:
        return [[2, 3], [1, 4]]
    mid = n // 2
    pre = solve(mid)
    res = [[0] * n for _ in range(n)]
    counter = [0, n * n // 4, n * n // 2, 3 * n * n // 4]
    for i in range(mid):
        for j in range(mid):
            res[i][j] = pre[i][j] + counter[1]
            res[mid + i][j] = pre[mid - 1 - j][mid - 1 - i] + counter[0]
            res[i][mid + j] = pre[i][j] + counter[2]
            res[mid + i][mid + j] = pre[j][i] + counter[3]
    return res


n = int(input())
res = solve(n)
for i in range(n):
    print(res[i][0], end="")
    for j in range(1, n):
        print(" %d" % res[i][j], end="")
    print()

# 输入：4
# 输出：
#      6 7 10 11
#      5 8 9 12
#      4 3 14 13
#      1 2 15 16
