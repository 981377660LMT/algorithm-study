from typing import List


MOD = int(1e9 + 7)

# 1 <= product.length <= 10 ^ 5
# 1 <= limit <= 10^9

# 如果一个一个投资，我们有一个贪心的策略：
# 每次选择价格最高的项目投资。这可以用堆来模拟，
# 但是本题 limit 高达 10^9 ，模拟是无法在时限内通过的。(这就和中华)
# 类似于 2141. 同时运行 N 台电脑的最长时间 1 <= batteries[i] <= 10^9 pq模拟会TLE
# 所以必须二分最后一次投资的价格 !


# 如果这道题二分想的是找到最小的价格，使得count<=limit 实际上这样做是不行的，因为最后剩下的不足minPrice
# 必须要找到最小的价格使得count>=limit，最后再减去超出的个数*minPrice即可

# 1648. 销售价值减少的颜色球
# https://leetcode-cn.com/problems/sell-diminishing-valued-colored-balls/
class Solution:
    def maxInvestment(self, product: List[int], limit: int) -> int:
        def check(mid: int) -> bool:
            res = 0
            for num in product:
                res += max(0, num - mid + 1)
                if res >= limit:
                    return True
            return res >= limit

        # 特判，恰好limit
        left, right = 0, int(1e10)
        while left <= right:
            mid = (left + right) >> 1
            """找到最小的价格使得count>=limit"""
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        if left == 0:
            return sum(p * (p + 1) // 2 for p in product) % MOD

        # 边界特判
        allCount = sum((p - left + 1) for p in product if p >= left)
        if allCount < limit:
            left -= 1

        res = 0
        count = 0
        for p in product:
            if p >= left:
                res += (p + left) * (p - left + 1) // 2
                count += p - left + 1

        res -= (count - limit) * left
        return res % MOD


print(Solution().maxInvestment(product=[4, 5, 3], limit=8))
print(Solution().maxInvestment(product=[2, 1, 3], limit=20))
print(Solution().maxInvestment(product=[3, 1, 2], limit=7))
print(
    Solution().maxInvestment(
        product=[43877, 10848, 10442, 48132, 83395, 71523, 60275, 39527], limit=345056
    )
)

