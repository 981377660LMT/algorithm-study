积分图 (Summed-area table)
积分图，又称总和面积表，是一个快速且有效的对一个网格的矩形子区域中计算和的数据结构和算法。

HashMap 是一个最常用的数据结构，它主要用于我们有通过固定值(key)获取内容的场景，时间复杂度可以最快优化到 O(1) 哈，当然效果不好的时候时间复杂度是 O(logN) 或者 O(n)。虽然固定值查找提高了速度，但是 HashMap 不能保证固定值，也就是 key 的顺序，所以这个时候 TreeMap 就出现了，虽然它的查找、删除、更新的时间复杂度都是 O(logN)，但是他可以保证 key 的有序性。
HashSet 使用的是 HashMap 来实现，而 TreeSet 使用的是 TreeMap 来实现的。

`前缀和`：p[i-1][j]+p[i][j-1]-p[i-1][j-1]+mat[i-1][j-1]
`矩形和`：p[x2][y2] - p[x1-1][y2] - p[x2][y1-1] + p[x1-1][y1-1]
`1292. 元素和小于等于阈值的正方形的最大边长.py`

```Python
m, n = len(mat), len(mat[0])
preSum = [[0] * (n + 1) for _ in range(m + 1)]
for r in range(1, m + 1):
    for c in range(1, n + 1):
        preSum[r][c] = (
            preSum[r - 1][c] + preSum[r][c - 1] - preSum[r - 1][c - 1] + mat[r - 1][c - 1]
        )

# 注意顺序x1,y1,x2,y2
def getSum(x1, y1, x2, y2):
    return preSum[x2][y2] - preSum[x1 - 1][y2] - preSum[x2][y1 - 1] + preSum[x1 - 1][y1 - 1]

```
