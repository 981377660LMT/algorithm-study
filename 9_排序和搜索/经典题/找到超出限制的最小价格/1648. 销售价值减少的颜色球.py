from typing import List


# 你有一些球的库存 inventory ，里面包含着不同颜色的球。一个顾客想要 任意颜色 总数为 orders 的球。
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

# 二分最后的同色球数量


class Solution:
    def maxProfit(self, inventory: List[int], orders: int) -> int:
        left = 0
        right = max(inventory)
        while left <= right:
            mid = (left + right) >> 1
            count = sum((inv - mid) for inv in inventory if inv > mid)
            # 超出orders
            if count >= orders:
                left = mid + 1
            else:
                right = mid - 1

        minPrice = left
        res, count = 0, 0
        for inv in inventory:
            if inv > minPrice:
                cur = inv - minPrice
                count += cur
                res += (inv + minPrice + 1) * cur // 2
                res %= MOD

        # 缺多少个，就补不多少个min_get_sum
        res += (orders - count) * minPrice
        res %= MOD
        return res


print(Solution().maxProfit(inventory=[2, 5], orders=4))
# 输出：14
# 解释：卖 1 个第一种颜色的球（价值为 2 )，卖 3 个第二种颜色的球（价值为 5 + 4 + 3）。
# 最大总和为 2 + 5 + 4 + 3 = 14 。

