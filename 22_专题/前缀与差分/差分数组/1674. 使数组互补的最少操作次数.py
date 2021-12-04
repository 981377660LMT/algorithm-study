# 每一次操作，你可以将 nums 中的任何整数替换为 1 到 limit 之间的另一个整数。
# 如果对于所有下标 i（下标从 0 开始），nums[i] + nums[n - 1 - i] 都等于同一个数，则数组 nums 是 互补的 。
# 返回使数组 互补 的 最少 操作次数。
# n 是偶数。2 <= n <= 105

# https://leetcode-cn.com/problems/minimum-moves-to-make-array-complementary/solution/jie-zhe-ge-wen-ti-xue-xi-yi-xia-chai-fen-shu-zu-on/
# 差分统计 以目标和为主线
# 我们考虑任意一个数对(a,b)，不妨假设a≤b。假设最终选定的和值为target
# 令target从数轴最左端开始向右移动 (1+a) (a+b) (a+b+1) (b+limit+1) 四个位置需要更新差分数组
# 最后，我们遍历（扫描）差分数组，就可以找到令总操作次数最小的target，同时可以得到最少的操作次数。
from typing import List


class Solution:
    def minMoves(self, nums: List[int], limit: int) -> int:
        diff = [0] * (2 * limit + 2)  # 差分数组
        n = len(nums)

        for i in range(n // 2):
            a, b = sorted((nums[i], nums[~i]))
            diff[a + 1] -= 1
            diff[a + b] -= 1
            diff[a + b + 1] += 1
            diff[b + limit + 1] += 1

        cur = n
        res = n
        for i in range(2, 2 * limit + 1):
            cur += diff[i]
            res = min(res, cur)
        return res


print(Solution().minMoves(nums=[1, 2, 4, 3], limit=4))
# 输出：1
# 解释：经过 1 次操作，你可以将数组 nums 变成 [1,2,2,3]（加粗元素是变更的数字）：
# nums[0] + nums[3] = 1 + 3 = 4.
# nums[1] + nums[2] = 2 + 2 = 4.
# nums[2] + nums[1] = 2 + 2 = 4.
# nums[3] + nums[0] = 3 + 1 = 4.
# 对于每个 i ，nums[i] + nums[n-1-i] = 4 ，所以 nums 是互补的。

