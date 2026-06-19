# 3950. 恰好一对连续置位
# https://leetcode.cn/problems/exactly-one-consecutive-set-bits-pair/description/
# 如果其二进制表示中 恰好 仅包含 一对 连续的置位 ，则返回 true ，否则返回 false 。


class Solution:
    def consecutiveSetBits(self, n: int) -> bool:
        m = n & (n >> 1)  # 所有相邻比特位的 &
        return m > 0 and m & (m - 1) == 0
