from typing import List

# 第一栋建筑的高度 必须 是 0 。
# 任意两栋相邻建筑的高度差 不能超过  1 。
# !某些建筑还有额外的最高高度限制
# 请你返回 最高 建筑能达到的 最高高度 。
# 建筑在数组 restrictions 中最多只会出现一次

# n<=10^9  => 说明是不可以遍历n的,遍历都超时
# 0 <= restrictions.length <= min(n - 1, 10^5) => 说明只应该遍历限制

# 1. 每一个限制实际上是对所有 n 栋建筑的限制
# 2. 需要从左到右/从右到左更新限制
# 3. 被限制的建筑高度应该等于限制(贪心)
# !4. i,j限制间建筑物的高度应该是一个山脉数组
# peek-left+peek-right<=j-i
# 临界有 peek=(left+right+j-i)//2


class Solution:
    def maxBuilding(self, n: int, rst: List[List[int]]) -> int:
        rst.extend([[1, 0], [n, n - 1]])
        rst = sorted(rst)

        # 1. 限制传递
        for i in range(len(rst) - 2, -1, -1):
            rst[i][1] = min(rst[i][1], rst[i + 1][1] + rst[i + 1][0] - rst[i][0])
        for i in range(1, len(rst)):
            rst[i][1] = min(rst[i][1], rst[i - 1][1] + rst[i][0] - rst[i - 1][0])

        # 2. 更新peek
        peek = 0
        for i in range(1, len(rst)):
            peek = max(peek, (rst[i - 1][1] + rst[i][1] + rst[i][0] - rst[i - 1][0]) // 2)
        return peek


print(Solution().maxBuilding(n=5, rst=[[2, 1], [4, 1]]))

# 输出：2
# 解释：上图中的绿色区域为每栋建筑被允许的最高高度。
# 我们可以使建筑高度分别为 [0,1,2,1,2] ，最高建筑的高度为 2 。
