from typing import List

# 每次查询需要找到一个非负整数 k < 2^maximumBit ，
# 使得 nums[0] XOR nums[1] XOR ... XOR nums[nums.length-1] XOR k 的结果 最大化
# 每次查询完后从当前数组 nums 删除 最后 一个元素。
# nums​​​ 中的数字已经按 升序 排好序。
# 0 <= nums[i] < 2^maximumBit

# 每次贪心即可 让最后异或结果为 111111...
class Solution:
    def getMaximumXor(self, nums: List[int], maximumBit: int) -> List[int]:
        maxXor = (1 << maximumBit) - 1
        res = []
        curXor = 0
        for num in nums:
            curXor ^= num
            res.append(maxXor ^ curXor)

        return res[::-1]


print(Solution().getMaximumXor(nums=[0, 1, 1, 3], maximumBit=2))
# 输出：[0,3,2,3]
# 解释：查询的答案如下：
# 第一个查询：nums = [0,1,1,3]，k = 0，因为 0 XOR 1 XOR 1 XOR 3 XOR 0 = 3 。
# 第二个查询：nums = [0,1,1]，k = 3，因为 0 XOR 1 XOR 1 XOR 3 = 3 。
# 第三个查询：nums = [0,1]，k = 2，因为 0 XOR 1 XOR 2 = 3 。
# 第四个查询：nums = [0]，k = 3，因为 0 XOR 3 = 3 。
