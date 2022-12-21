from typing import List

# 你的开销是单个袋子里球数目的 最大值 ，你想要 最小化 开销。
# 1 <= nums.length <= 1e5


class Solution:
    def minimumSize(self, nums: List[int], maxOperations: int) -> int:
        def check(mid: int) -> bool:
            """拆分后每个数字都不超过mid"""
            count = 0
            for num in nums:
                # divmod均等拆成div份
                div, mod = num // mid, num % mid
                count += div - 1 if mod == 0 else div
                if count > maxOperations:
                    return False
            return True

        left, right = 1, int(1e16)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(Solution().minimumSize(nums=[2, 4, 8, 2], maxOperations=4))
# 输出：2
# 解释：
# - 将装有 8 个球的袋子分成装有 4 个和 4 个球的袋子。[2,4,8,2] -> [2,4,4,4,2] 。
# - 将装有 4 个球的袋子分成装有 2 个和 2 个球的袋子。[2,4,4,4,2] -> [2,2,2,4,4,2] 。
# - 将装有 4 个球的袋子分成装有 2 个和 2 个球的袋子。[2,2,2,4,4,2] -> [2,2,2,2,2,4,2] 。
# - 将装有 4 个球的袋子分成装有 2 个和 2 个球的袋子。[2,2,2,2,2,4,2] -> [2,2,2,2,2,2,2,2] 。
# 装有最多球的袋子里装有 2 个球，所以开销为 2 并返回 2 。
