from typing import List

# 1 <= nums.length <= 105

# 402. 移掉 K 位数字.py
class Solution:
    def mostCompetitive(self, nums: List[int], k: int) -> List[int]:
        remove = len(nums) - k
        stack = []
        for char in nums:
            while remove and stack and stack[-1] > char:
                stack.pop()
                remove -= 1
            stack.append(char)
        return stack[:k]


print(Solution().mostCompetitive(nums=[3, 5, 2, 6], k=2))
# 输出：[2,6]
# 解释：在所有可能的子序列集合 {[3,5], [3,2], [3,6], [5,2], [5,6], [2,6]} 中，[2,6] 最具竞争力。
