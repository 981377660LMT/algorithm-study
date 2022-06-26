from typing import List

# https://leetcode.cn/problems/largest-combination-with-bitwise-and-greater-than-zero/
# 2275. 按位与结果大于零的最长组合


class Solution:
    def largestCombination(self, candidates: List[int]) -> int:
        """
        返回按位与结果大于 0 的 最长组合(子序列)的长度。

        统计每位是否为1
        时间复杂度 O(nlogA)
        """
        bitCount = [0] * 32
        for num in candidates:
            for i in range(32):
                if num & (1 << i):
                    bitCount[i] += 1

        return max(bitCount)

