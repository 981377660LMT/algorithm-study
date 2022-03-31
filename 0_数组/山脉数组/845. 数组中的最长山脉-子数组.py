from typing import List

# 数组 A 中符合下列属性的任意`连续子数组` B 称为 “山脉”：
# B.length >= 3
# 先绝对单增后绝对单减

# 返回最长 “山脉” 的长度。


class Solution:
    def longestMountain(self, arr: List[int]) -> int:
        up, down = [0] * len(arr), [0] * len(arr)
        for i in range(1, len(arr)):
            if arr[i] > arr[i - 1]:
                up[i] = up[i - 1] + 1
        for i in range(len(arr) - 2, -1, -1):
            if arr[i] > arr[i + 1]:
                down[i] = down[i + 1] + 1
        return max([u + d + 1 for u, d in zip(up, down) if u and d], default=0)


print(Solution().longestMountain([2, 1, 4, 7, 3, 2, 5]))
# 输入：[2,1,4,7,3,2,5]
# 输出：5
# 解释：最长的 “山脉” 是 [1,4,7,3,2]，长度为 5。
