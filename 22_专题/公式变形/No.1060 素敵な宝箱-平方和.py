# 有n个宝箱，每个宝箱里有若干种宝石
# Alice和Bob依次选取宝箱，每次只能选取一个宝箱，选取后宝箱里的宝石全部归自己所有
# 得分为C1^2+C2^2+...+Cn^2,其中Ci为第i种宝石的数量
# !求Alice-Bob的得分差值最大值

# 公式变形：
# 对第i种宝石，假设有ci个
# !Alice有x个，那么差值贡献为 x*x-(ci-x)*(ci-x) = ci*x-ci*(ci-x)
# !可以理解为第i种宝石每个ci分

from typing import List


def yuki1060(grid: List[List[int]]) -> int:
    ROW, COL = len(grid), len(grid[0])
    colSum = [0] * COL
    for i in range(ROW):
        for j in range(COL):
            colSum[j] += grid[i][j]
    rowPrice = [0] * ROW
    for i in range(ROW):
        for j in range(COL):
            rowPrice[i] += grid[i][j] * colSum[j]
    rowPrice.sort(reverse=True)
    res1 = sum(rowPrice[::2])
    res2 = sum(rowPrice[1::2])
    return res1 - res2


if __name__ == "__main__":
    n, m = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(n)]
    print(yuki1060(grid))
