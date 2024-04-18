# https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/solutions/2511644/tu-jie-jian-ji-gao-xiao-yi-tu-miao-dong-1427d/
# 试填法
# 最高位能不能是1？->次高位能不能是1？->...
# 为了找到两个数异或等于11000，哈希表保存即可.
# 将最大化问题变成判定性问题


from typing import List


class Solution:
    def findMaximumXOR(self, nums: List[int]) -> int:
        maxs = max(nums, default=0)
        bitLen = maxs.bit_length()
        res, mask = 0, 0
        for b in range(bitLen - 1, -1, -1):
            mask |= 1 << b
            cand = res | (1 << b)  # 试填，这个比特位可以是1吗？
            visited = set()
            for v in nums:
                v &= mask  # 低于b位的置为0
                if cand ^ v in visited:
                    res = cand
                    break
                visited.add(v)
        return res
