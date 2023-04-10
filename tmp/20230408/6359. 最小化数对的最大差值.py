# 给你一个下标从 0 开始的整数数组 nums 和一个整数 p 。
# !请你从 nums 中找到 p 个下标对，每个下标对对应数值取差值，
# 你需要使得这 p 个差值的 最大值 最小。
# 同时，你需要确保每个下标在这 p 个下标对中最多出现一次。

# 对于一个下标对 i 和 j ，这一对的差值为 |nums[i] - nums[j]| ，
# 其中 |x| 表示 x 的 绝对值 。
# 请你返回 p 个下标对对应数值 最大差值 的 最小值 。

from typing import List


class Solution:
    def minimizeMax(self, nums: List[int], p: int) -> int:
        if p == 0:
            return 0
        nums.sort()

        def check(mid: int) -> bool:
            """选p对使得每一对差值不超过mid"""
            # !相邻选最好
            visited = [False] * n
            ok = 0
            for i in range(n - 1):
                if visited[i]:
                    continue
                if nums[i + 1] - nums[i] <= mid:
                    ok += 1
                    visited[i] = True
                    visited[i + 1] = True
                    if ok >= p:
                        return True
            return ok >= p

        n = len(nums)
        left, right = 0, int(1e12)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
