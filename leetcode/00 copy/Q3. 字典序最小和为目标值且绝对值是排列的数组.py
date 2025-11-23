# 给你一个正整数 n 和一个整数 target。
#
# 请返回一个大小为 n 的 字典序最小 的整数数组，并满足：
#
# 其元素 和 等于 target。
# 其元素的 绝对值 组成一个大小为 n 的 排列。
# 如果不存在这样的数组，则返回一个空数组。
#
# 如果数组 a 和 b 在第一个不同的位置上，数组 a 的元素小于 b 的对应元素，则认为数组 a 字典序小于 数组 b。
#
# 大小为 n 的 排列 是对整数 1, 2, ..., n 的重新排列

from typing import List


class Solution:
    def lexSmallestNegatedPerm(self, n: int, target: int) -> List[int]:
        allSum = n * (n + 1) // 2
        if (allSum - target) % 2 == 1 or abs(target) > allSum:
            return []
        remain = (allSum - target) // 2
        neg = [False] * (n + 1)
        cand = n
        while remain > 0 and cand > 0:
            if remain >= cand:
                neg[cand] = True
                remain -= cand
            cand -= 1
        if remain > 0:
            return []
        res1 = [-i for i in range(n, 0, -1) if neg[i]]
        res2 = [i for i in range(1, n + 1) if not neg[i]]
        return res1 + res2
