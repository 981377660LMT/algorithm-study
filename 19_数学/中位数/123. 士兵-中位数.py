# 士兵可以进行移动，每次移动，一名士兵可以向上，向下，向左或向右移动一个单位
# （因此，他的 x 或 y 坐标也将加 1 或减 1）。
# 现在希望通过移动士兵，使得所有士兵彼此相邻的处于同一条水平线内，
# 即所有士兵的 y 坐标相同并且 x 坐标相邻。
# 请你计算满足要求的情况下，所有士兵的总移动次数最少是多少。
# 需注意，两个或多个士兵不能占据同一个位置。


# 1. 上下移动与左右移动可以分开进行
# 2. 找出y坐标的中位数，用最小的代价先将所有点移到同一列上
# 3. 将这些点彼此相邻


n = int(input())
xPos = []
yPos = []
for _ in range(n):
    x, y = map(int, input().split())
    xPos.append(x)
    yPos.append(y)


xPos.sort()
yPos.sort()

res = 0
# 移到同一列上
mid = yPos[n // 2]
res += sum(abs(y - mid) for y in yPos)

# 所有数移到相邻
# 1703. 得到连续 K 个 1 的最少相邻交换次数
# abs(x0 - a) + abs(x1 - 1 - a) + abs(x2 - 2 - a) + … + abs(x(n-1) - (n - 1) - a)
xPos2 = [x - i for i, x in enumerate(xPos)]
xPos2.sort()
mid = xPos2[n // 2]
res += sum(abs(x - mid) for x in xPos2)

print(res)

