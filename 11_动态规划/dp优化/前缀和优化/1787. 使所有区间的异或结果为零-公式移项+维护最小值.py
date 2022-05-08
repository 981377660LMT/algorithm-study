# 1 <= k <= nums.length <= 2000
# ​​​​​​0 <= nums[i] < 2^10
# 返回数组中 要更改的最小元素数 ，以使所有长度为 k 的区间异或结果等于零。
# 即如何决策前k个数值，让最终的数组和原始数组的差异个数最少
from typing import List
from collections import defaultdict


class Solution:
    def minChanges(self, nums: List[int], k: int) -> int:
        """画成二维矩阵比较好思考
        
        列间转移：暴力尝试1024种肯定会超时(k*2^20)
        删边优化；优化为仅修改这一列的一部分数/修改这一列的全部数
        注意到转移式子dp[col][xor]=min(dp[col-1][xor^select]-counters[col][select])+counts[col]
        如果select没在列中出现 即为 min(dp[col-1][xor^select]) 对此只需用dp[col-1]最小值即可
        如果select在列中出现 即为 min(dp[col-1][xor^select]-counters[col][select])
        从而大大减少了列间暴力转移的边数
        观察转移式优化dp
        """
        counters = [defaultdict(int) for _ in range(k)]  # 用counter超时了 改defaultdict
        counts = [0] * k
        for i, num in enumerate(nums):
            col = i % k
            counters[col][num] += 1
            counts[col] += 1

        # d[col][xor] 且首行前col列异或值为xor时，最少需要更改的次数
        dp = [int(1e20)] * 1024  # 优化成一维数组 dp 和 ndp
        dp[0] = 0
        for i in range(1, k + 1):
            preMin = min(dp)
            ndp = [preMin + counts[i - 1]] * 1024
            for pre in range(1024):
                for key in counters[i - 1].keys():
                    cur = pre ^ key
                    ndp[cur] = min(ndp[cur], dp[pre] + counts[i - 1] - counters[i - 1][key])
            dp = ndp

        return dp[0]


# 3 3
print(Solution().minChanges(nums=[3, 4, 5, 2, 1, 7, 3, 4, 7], k=3))
print(Solution().minChanges(nums=[1, 2, 4, 1, 2, 5, 1, 2, 6], k=3))

# 输出：3
# 解释：将数组 [3,4,5,2,1,7,3,4,7] 修改为 [3,4,7,3,4,7,3,4,7]

