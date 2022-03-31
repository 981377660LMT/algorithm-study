# O(1)空间复杂度检查是否所有数freq为偶数

# 异或,但是注意1^2^3为0
# 异或之前先哈希
class Solution:
    def solve(self, nums):
        xor_ = 0
        for num in nums:
            hash = ((num + 0x9E3779B1) * 0x85EBCA77) & 0xFFFFFFFF
            xor_ ^= hash
        return xor_ == 0

