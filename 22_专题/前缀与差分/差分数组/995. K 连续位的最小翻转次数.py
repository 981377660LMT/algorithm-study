from collections import defaultdict
from typing import List

# !1871. 跳跃游戏-差分范围更新
# k位翻转 就是从 nums 中选择一个长度为 k 的 子数组 ，
# 同时把子数组中的每一个 0 都改成 1 ，把子数组中的每一个 1 都改成 0 。
# 返回数组中`不存在 0` 所需的最小 k位翻转 次数。如果不可能，则返回 -1 。
# k个一组翻转


class Solution:
    def minKBitFlips(self, nums: List[int], k: int) -> int:
        """返回数组中不存在 0 所需的最小 k位翻转 次数;如果不可能，则返回 -1 。"""
        diff = defaultdict(int)  # 记录翻转次数
        res = 0
        for i, num in enumerate(nums):
            diff[i] += diff[i - 1]
            if (diff[i] & 1) ^ num == 0:  # 需要反转了
                if i + k - 1 >= len(nums):
                    return -1
                diff[i] += 1
                diff[(i + k - 1) + 1] -= 1
                res += 1
        return res


print(Solution().minKBitFlips([0, 0, 0, 1, 0, 1, 1, 0], 3))
