# 2952. 需要添加的硬币的最小数量-二进制分桶
# https://leetcode.cn/problems/minimum-number-of-coins-to-be-added/solutions/3052162/xian-xing-suan-fa-by-hqztrue-rf5n/


INF = int(1e18)


class Solution:
    def minimumAddedCoins(self, coins: list[int], target: int) -> int:
        max_ = max(coins)
        bitLen = max_.bit_length()
        bucketMin = [INF] * (bitLen + 1)
        bucketSum = [0] * (bitLen + 1)

        res = 0
        s = 0

        for x in coins:
            b = x.bit_length()
            if bucketSum[b] < target:
                bucketSum[b] += x
            bucketMin[b] = min(bucketMin[b], x)

        for min_, sum_ in zip(bucketMin, bucketSum):
            if min_ != INF:
                while s < target and s + 1 < min_:
                    res += 1
                    s += s + 1
            s += sum_

        while s < target:
            res += 1
            s += s + 1

        return res
