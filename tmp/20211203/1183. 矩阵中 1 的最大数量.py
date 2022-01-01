from heapq import heappop, heappush

# 矩阵 M 中每个大小为 sideLength * sideLength 的 正方形 子阵中，1 的数量不得超过 maxOnes。
# 请你设计一个算法，计算矩阵中最多可以有多少个 1。

# 假设在某个点(i, j)上放置一个1，则可以再所有满足x % sideLength == i和y % sideLength == j的位置(x, y)放置一个1且互相之间不影响
# 所以只需要找出在第一个边长为sideLength正方形内的哪些位置放置1能使得整个矩形内的1最多即可
# 遍历第一个边长为sideLength正方形内的每个点，找出前maxOnes个能使得在矩阵内放尽可能多的点即可。

# 计算左上角正方形的每个格子在整个矩形中有多少个等效位置，取等效位置最多的前maxOnes个即可
class Solution:
    def maximumNumberOfOnes(self, width: int, height: int, sideLength: int, maxOnes: int) -> int:
        pq = []

        for i in range(sideLength):
            for j in range(sideLength):
                rowCount = ((width - i - 1) // sideLength) + 1
                colCount = ((height - j - 1) // sideLength) + 1
                heappush(pq, rowCount * colCount)
                if len(pq) > maxOnes:
                    heappop(pq)

        return sum(pq)


print(Solution().maximumNumberOfOnes(width=3, height=3, sideLength=2, maxOnes=1))
# 输出：4
# 解释：
# 题目要求：在一个 3*3 的矩阵中，每一个 2*2 的子阵中的 1 的数目不超过 1 个。
# 最好的解决方案中，矩阵 M 里最多可以有 4 个 1，如下所示：
# [1,0,1]
# [0,0,0]
# [1,0,1]
