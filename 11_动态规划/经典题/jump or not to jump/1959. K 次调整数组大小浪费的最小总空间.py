# 你正在设计一个动态数组
# 其中 nums[i] 是 i 时刻数组中的元素数目
# k表示你可以 调整 数组大小的 最多 次数(k表示你可以 调整 数组大小的 最多 次数)
# 在调整数组大小不超过 k 次的前提下，请你返回 最小总浪费空间 。
# 注意：数组最开始时可以为 任意大小 ，且 不计入 调整大小的操作次数。
# 1 <= nums.length <= 200

# 总结：我们可以变换k次，等价于我们可以`将原数组分成k+1个区间，每个区间必须取区间最大值`
# 每个区间都要求和（区间最大值*区间长度再减去区间和，）
from typing import List
from functools import lru_cache
from itertools import accumulate

INF = 0x3FFFFFFF


class Solution:
    def minSpaceWastedKResizing(self, nums: List[int], k: int) -> int:
        preSum = [0] + list(accumulate(nums))

        @lru_cache(None)
        def dfs(cur: int, remain: int) -> int:
            if cur == len(nums):
                return 0
            if remain < 0:
                return INF

            res = INF

            curMax = 0
            for next in range(cur, len(nums) - remain):
                curMax = max(curMax, nums[next])
                res = min(
                    res,
                    curMax * (next - cur + 1)
                    - (preSum[next + 1] - preSum[cur])
                    + dfs(next + 1, remain - 1),
                )

            return res

        return dfs(0, k)


print(Solution().minSpaceWastedKResizing(nums=[10, 20, 30], k=1))
# 输出：10
# 解释：size = [20,20,30].
# 我们可以让数组初始大小为 20 ，然后时刻 2 调整大小为 30 。
# 总浪费空间为 (20 - 10) + (20 - 20) + (30 - 30) = 10 。

