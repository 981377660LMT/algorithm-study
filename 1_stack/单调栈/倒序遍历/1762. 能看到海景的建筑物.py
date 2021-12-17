from typing import List


class Solution:
    def findBuildings(self, heights: List[int]) -> List[int]:
        invalidIndex = set()
        stack = []

        for i in range(len(heights) - 1, -1, -1):
            # 如果进来的矮的，则什么都不用做
            # 只处理进来矮的情况
            if not stack or stack[-1] < heights[i]:
                stack.append(heights[i])
            else:
                invalidIndex.add(i)

        return [i for i in range(len(heights)) if i not in invalidIndex]


print(Solution().findBuildings([4, 2, 3, 1]))
# 输出：[0,2,3]
# 解释：1 号建筑物看不到海景，因为 2 号建筑物比它高
