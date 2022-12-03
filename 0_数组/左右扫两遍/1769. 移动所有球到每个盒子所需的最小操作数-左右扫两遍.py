# 输入：boxes = "110"
# 输出：[1,1,3]
# 解释：每个盒子对应的最小操作数如下：
# 1) 第 1 个盒子：将一个小球从第 2 个盒子移动到第 1 个盒子，需要 1 步操作。
# 2) 第 2 个盒子：将一个小球从第 1 个盒子移动到第 2 个盒子，需要 1 步操作。
# 3) 第 3 个盒子：将一个小球从第 1 个盒子移动到第 3 个盒子，需要 2 步操作。将一个小球从第 2 个盒子移动到第 3 个盒子，需要 1 步操作。共计 3 步操作。

# 其中 answer[i] 是将所有小球移动到第 i 个盒子所需的 最小 操作数。
# 每次只能相邻移动1个球
# 1 <= n <= 2000

# !前后缀分解 移动总数=前缀移动数+后缀移动数
from typing import List


class Solution:
    def minOperations(self, boxes: str) -> List[int]:
        def makeDp(nums: List[int]) -> List[int]:
            n = len(nums)
            dp = [0] * (n + 1)
            count = 0
            for i in range(n):
                count += nums[i]
                dp[i + 1] = count + dp[i]
            return dp

        nums = list(map(int, boxes))
        n = len(nums)
        preDp, sufDp = makeDp(nums), makeDp(nums[::-1])[::-1]
        return [preDp[i] + sufDp[i + 1] for i in range(n)]


print(Solution().minOperations("110"))
