from typing import List

# 数组 A 中符合下列属性的任意`子序列` B 称为 “山脉”：
# B.length >= 3
# 先绝对单增后绝对单减

# 返回最长 “山脉” 的长度。
# 数组中求一个最长的山脉的长度，满足山脉中没有两个相邻的相同高度的点，且只有一个峰

# n<=1000


class Solution:
    def longestMountain(self, arr: List[int]) -> int:
        n = len(arr)
        up = [0] * n
        down = [0] * n

        # 前后预处理+查表
        for i in range(n):
            for j in range(i):
                if arr[i] > arr[j]:
                    up[i] = max(up[i], up[j] + 1)

        arr = arr[::-1]
        for i in range(n):
            for j in range(i):
                if arr[i] > arr[j]:
                    down[i] = max(down[i], down[j] + 1)
        down = down[::-1]

        return max([u + d + 1 for u, d in zip(up, down)], default=1)


print(Solution().longestMountain([2, 1, 4, 7, 3, 2, 5]))
# 输入：[2,1,4,7,3,2,5]
# 输出：5
# 解释：最长的 “山脉” 是 [1,4,7,3,2]，长度为 5。
