from typing import List

# 返回数组中 要更改的最小元素数 ，以使所有长度为 k 的区间异或结果等于零。

# One of the hardest
class Solution:
    def minChanges(self, nums: List[int], k: int) -> int:
        ...


print(Solution().minChanges(nums=[1, 2, 4, 1, 2, 5, 1, 2, 6], k=3))
# 输出：3
# 解释：将数组[1,2,4,1,2,5,1,2,6] 修改为 [1,2,3,1,2,3,1,2,3]


# 直接放弃
