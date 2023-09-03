# https://leetcode.cn/problems/minimum-operations-to-form-subsequence-with-target-sum/
#
# 给你一个下标从 0 开始的数组 nums ，它包含 非负 整数，且全部为 2 的幂，同时给你一个整数 target 。
# 一次操作中，你必须对数组做以下修改：
# 选择数组中一个元素 nums[i] ，满足 nums[i] > 1 。
# 将 nums[i] 从数组中删除。
# 在 nums 的 末尾 添加 两个 数，值都为 nums[i] / 2 。
# !你的目标是让 nums 的一个 子序列 的元素和等于 target ，请你返回达成这一目标的 最少操作次数 。
# 如果无法得到这样的子序列，请你返回 -1 。
# nums 只包含非负整数，且均为 2 的幂。
#
# n<=1e5,nums[i]<=2**30,1<=target<2**31
#
# https://leetcode.cn/problems/minimum-operations-to-form-subsequence-with-target-sum/solutions/2413344/tan-xin-by-endlesscheng-immn/
# 从低位到高位贪心：
# - 如果target第i位为0，跳过
# - 如果target第i为为1，看<=2**i的元素和能否>=target&((1<<(i+1))-1)
#   !- 如果能，那么ok (因为可以凑出2*0,2*1,2*2,...,2**(i-1)且和大于等于target&((1<<(i+1))-1)
#   - 如果不能，那么就需要把一个更大的数 (2**j)不断一份为二直到分解出2**i，同时无需判断(i+1到j-1)位


from typing import Counter, List


class Solution:
    def minOperations(self, nums: List[int], target: int) -> int:
        if sum(nums) < target:
            return -1
        counter = Counter(nums)
        res = 0
        curSum = 0
        bit = 0
        while 1 << bit <= target:
            curSum += counter[1 << bit] << bit
            mask = (1 << (bit + 1)) - 1
            bit += 1
            if curSum < target & mask:
                res += 1
                while counter[1 << bit] == 0:
                    res += 1
                    bit += 1
        return res
