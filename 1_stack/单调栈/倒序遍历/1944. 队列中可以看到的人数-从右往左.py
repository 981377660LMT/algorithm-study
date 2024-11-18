# 1944. 队列中可以看到的人数
# https://leetcode.cn/problems/number-of-visible-people-in-a-queue/description/
# 请你返回一个长度为 n 的数组 answer ，其中 answer[i] 是第 i 个人在他右侧队列中能 看到 的 人数 。
# 给你以一个整数数组 heights ，每个整数 `互不相同`，heights[i] 表示第 i 个人的高度。

from typing import List


class Solution:
    def canSeePersonsCount(self, heights: List[int]) -> List[int]:
        res = [0] * len(heights)
        stack = []  # 单调递减栈
        for i in range(len(heights) - 1, -1, -1):
            # 如果进来的矮的，则什么都不用做
            # 只处理进来高的情况
            while stack and heights[stack[-1]] < heights[i]:
                res[i] += 1  # 看到一个人
                stack.pop()  # 前面那个对之后没影响了
            # 矮的前面有一个高的挡住了视线
            if stack:
                res[i] += 1
            while stack and heights[stack[-1]] == heights[i]:  # !相等的全部出栈，只能算一次
                stack.pop()
            stack.append(i)
        return res


print(Solution().canSeePersonsCount(heights=[10, 6, 8, 5, 11, 9]))
print(Solution().canSeePersonsCount(heights=[3, 1, 4, 2, 5]))
# 输出：[3,1,2,1,1,0]
# 解释：
# 第 0 个人能看到编号为 1 ，2 和 4 的人。
# 第 1 个人能看到编号为 2 的人。
# 第 2 个人能看到编号为 3 和 4 的人。
# 第 3 个人能看到编号为 4 的人。
# 第 4 个人能看到编号为 5 的人。
# 第 5 个人谁也看不到因为他右边没人。
