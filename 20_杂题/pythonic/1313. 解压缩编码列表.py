from typing import List


class Solution:
    def decompressRLElist(self, nums: List[int]) -> List[int]:
        return sum(([b] * a for a, b in zip(nums[::2], nums[1::2])), [])


# 输入：nums = [1,2,3,4]
# 输出：[2,4,4,4]
# 解释：第一对 [1,2] 代表着 2 的出现频次为 1，所以生成数组 [2]。
# 第二对 [3,4] 代表着 4 的出现频次为 3，所以生成数组 [4,4,4]。
# 最后将它们串联到一起 [2] + [4,4,4] = [2,4,4,4]。

print(sum([[1], [2]], []))

