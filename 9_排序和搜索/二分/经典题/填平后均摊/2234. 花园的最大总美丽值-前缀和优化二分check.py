from itertools import accumulate
from typing import List


# 2234. 花园的最大总美丽值-前缀和优化二分check
# 每次操作可以给数组中的某个数增加1,最多进行k次操作。
# 最后获得的价值为：数组中大于t的数的数量cnt乘以a 加上 小于t的数里面的最小的数minv乘以b。
# 问最大价值是多少？


# 代码能力过弱以至于没撸出纯码的T4
# 深有同感，知道怎么写，不过写着写着就写乱了，导致最后调试都不方便
# 非常乱的逻辑，应该要抽离成函数
# 第四题大思路是对的，不过二分的验证函数写的不合适，既然已经有了前缀/后缀和，
# 那单次验证就可以O(1)，而我写的单次验证是O(n)，总复杂度就成了O(n^2 logn)，肯定过不去了……

# !枚举最后一共有几座完善的花园，再二分不完善的花园中，花的最少数目最大可以是多少


class Solution:
    def maximumBeauty(
        self, flowers: List[int], k: int, target: int, full: int, partial: int
    ) -> int:
        flowers.sort()
        n = len(flowers)
        preSum = [0] + list(accumulate(flowers))
        nums = [0] + flowers  # 哨兵，与前缀和统一形式

        # 特判
        if flowers[0] >= target:
            return full * n

        moreThan = sum(f >= target for f in flowers)  # 已经满足的
        res = 0
        cost = 0  # 表示为了获得规定数量的完善花园，需要多种几朵花(已经用完的花的数量)

        # 枚举full的个数
        for i in range(moreThan, n + 1):
            if cost >= k:
                break

            # 最右二分，用剩下的花，可以把不完善的花园里`最少的花`变得和`哪个不完善的花园`一样多``，
            # 然后加上`剩下的除以这个数量`就是最小值，即填平后均摊
            # 所有花园里花的数目至少和第 head 个花园一样多，可能还剩下一些花，还能再提高花的数目
            # 不完善花园里的花不能超过 target，否则就变成完善花园了
            left, right = 0, n - i
            while left <= right:
                mid = (left + right) // 2
                diff = nums[mid] * mid - preSum[mid]  # 关键是自己这里没用前缀和，写的nlogn 用了前缀和就是O(logn)
                if cost + diff <= k:
                    left = mid + 1
                else:
                    right = mid - 1

            # 此时最少的花可以变为right
            remain = k - cost - (right * nums[right] - preSum[right])
            min_ = min(target - 1, nums[right] + remain // right if right else 0)
            res = max(res, full * i + min_ * partial)
            cost += target - nums[n - i]

        return res


print(Solution().maximumBeauty(flowers=[1, 3, 1, 1], k=7, target=6, full=12, partial=1))
