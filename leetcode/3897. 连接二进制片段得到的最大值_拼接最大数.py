# 3897. 连接二进制片段得到的最大值_拼接最大数
# https://leetcode.cn/problems/maximum-value-of-concatenated-binary-segments/description/
#
# 纯 1 片段放最前
# 纯 0 片段放最后
# 其余混合片段按 1 的数量降序排
# 若 1 的数量相同，按 0 的数量升序排

from typing import List

MOD = int(1e9 + 7)
MX = int(1e5 + 10)
pow2 = [1] * MX
for i in range(1, MX):
    pow2[i] = pow2[i - 1] * 2 % MOD


class Solution:
    def maxValue(self, nums1: List[int], nums0: List[int]) -> int:
        # 没有 0 的排在最前面，1 多的排前面，0 少的排前面
        order = sorted(range(len(nums1)), key=lambda i: (nums0[i] != 0, -nums1[i], nums0[i]))
        res = 0
        for i in order:
            res = ((res + 1) * pow2[nums1[i]] - 1) * pow2[nums0[i]] % MOD
        return res
