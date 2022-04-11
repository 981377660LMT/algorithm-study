from typing import List


# 你有一些球的库存 inventory ，里面包含着不同颜色的球。
# 一个顾客想要 任意颜色 总数为 orders 的球。
# 每个球的价值是目前剩下的 同色球 的数目。
# 请你返回卖了 orders 个球以后 最大 总价值之和

MOD = int(1e9 + 7)

# 1 <= inventory.length <= 105
# 1 <= inventory[i] <= 109

# 总结：
# 注意order可以很大 常规方法一次只pop/push 一个会超时
# 二分查找：
# 我们可以二分查找 `最后一次卖出时，球的价格 x`(最左二分)
# 最后的答案由两部分组成：
# 对于所有数量 > x 的颜色，其肯定会减小到 x，因此用等差数列求和公式求和即可。
# 如果执行完第 1 步，仍有剩余的 orders，则这些 orders 一定会以价格 x 卖出。
# https://leetcode-cn.com/problems/sell-diminishing-valued-colored-balls/solution/liang-chong-si-lu-you-hua-tan-xin-suan-fa-you-hua-/

# 注意这个最左二分需要大于等于orders


class Solution:
    def maxProfit(self, inventory: List[int], orders: int) -> int:
        def check(mid: int) -> bool:
            """订单数大于order的最后卖出最小价值"""
            count = 0
            for num in inventory:
                count += max(0, num - mid)
                if count > orders:
                    return True
            return False

        left = 0
        right = int(1e10)
        while left <= right:
            mid = (left + right) >> 1
            # 超出orders才回升
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        # # check边界
        # if sum((i - left + 1) for i in inventory if i >= left) < orders:
        #     left -= 1

        minPrice = left
        res, count = 0, 0
        for inv in inventory:
            if inv >= minPrice:
                cur = inv - minPrice + 1
                count += cur
                res += (minPrice + inv) * cur // 2
                res %= MOD

        # 超出多少个以minPrice卖出的
        res -= (count - orders) * minPrice
        res %= MOD
        return res


print(Solution().maxProfit(inventory=[2, 5], orders=4))
# 输出：14
# 解释：卖 1 个第一种颜色的球（价值为 2 )，卖 3 个第二种颜色的球（价值为 5 + 4 + 3）。
# 最大总和为 2 + 5 + 4 + 3 = 14 。

