from typing import List

# 1 <= n <= 15
# -104 <= sums[i] <= 104
class Solution:
    def recoverArray(self, n: int, sums: List[int]) -> List[int]:
        ...


print(Solution().recoverArray(n=4, sums=[0, 0, 5, 5, 4, -1, 4, 9, 9, -1, 4, 3, 4, 8, 3, 8]))
# 输出：[0,-1,4,5]
# 解释：[0,-1,4,5] 能够满足给出的子集的和。
