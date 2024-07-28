# 1674. 使数组互补的最少操作次数
# https://leetcode-cn.com/problems/minimum-moves-to-make-array-complementary/solution/jie-zhe-ge-wen-ti-xue-xi-yi-xia-chai-fen-shu-zu-on/
# 每一次操作，你可以将 nums 中的任何整数替换为 1 到 limit 之间的另一个整数。
# 如果对于所有下标 i（下标从 0 开始），nums[i] + nums[n - 1 - i] 都等于同一个数，则数组 nums 是 互补的 。
# 返回使数组 互补 的 最少 操作次数。
# 1 <= nums[i] <= limit <= 1e5
# n 是偶数
#
# 差分统计 `给 [l, r] 的区间加上一个数字 a, 只需要 diff[l] += a，diff[r + 1] -= a。`
# 我们考虑任意一个数对(a,b)，不妨假设a≤b。假设最终选定的和值为target
# 令target从数轴最左端开始向右移动 (1+a) (a+b) (a+b+1) (b+limit+1) 四个位置需要更新差分数组
# 最后，我们遍历（扫描）差分数组，就可以找到令总操作次数最小的target，同时可以得到最少的操作次数。
#
# 3224. 使差值相等的最少数组改动次数
# https://leetcode.cn/problems/minimum-array-changes-to-make-differences-equal/description/

from collections import defaultdict
from typing import List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minMoves(self, nums: List[int], limit: int) -> int:
        diff = defaultdict(int)  # 最终互补的数字和为 x，需要的操作数

        def add(left: int, right: int, x: int) -> None:
            diff[left] += x
            diff[right + 1] -= x

        n = len(nums)
        for i in range(n // 2):
            a, b = nums[i], nums[~i]
            if a > b:
                a, b = b, a

            add(2, a, 2)  #  2 <= x <= a 时需要两次操作
            add(a + 1, a + b - 1, 1)
            add(a + b + 1, b + limit, 1)
            add(b + limit + 1, INF, 2)

        res, curSum = INF, 0
        for key in sorted(diff):
            if 2 <= key <= 2 * limit:  # [2,2*limit]间的最小值
                curSum += diff[key]
                res = min2(res, curSum)
        return res
