from typing import List


class Solution:
    def canSeePersonsCount(self, heights: List[int]) -> List[int]:
        res = [0] * len(heights)
        stack = []
        for i in range(len(heights) - 1, -1, -1):
            # 如果进来的矮的，则什么都不用做
            # 只处理进来高的情况
            while stack and stack[-1] <= heights[i]:
                res[i] += 1
                stack.pop()  # 前面那个对之后没影响了

            # 矮的前面有一个高的挡住了视线
            if stack:
                res[i] += 1

            stack.append(heights[i])
        return res


print(Solution().canSeePersonsCount(heights=[10, 6, 8, 5, 11, 9]))
# 输出：[3,1,2,1,1,0]
# 解释：
# 第 0 个人能看到编号为 1 ，2 和 4 的人。
# 第 1 个人能看到编号为 2 的人。
# 第 2 个人能看到编号为 3 和 4 的人。
# 第 3 个人能看到编号为 4 的人。
# 第 4 个人能看到编号为 5 的人。
# 第 5 个人谁也看不到因为他右边没人。
