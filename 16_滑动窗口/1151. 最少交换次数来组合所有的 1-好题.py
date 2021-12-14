from typing import List

# 你需要通过交换位置，将数组中 任何位置 上的 1 组合到一起，并返回所有可能中所需 最少的交换次数。

# 固定区间，假设最后的全1就在当前的窗口
# 则窗口内0的个数==需要交换的次数


class Solution:
    def minSwaps(self, data: List[int]) -> int:
        n = len(data)
        windowSize = data.count(1)
        curOnes = sum(data[:windowSize])
        res = windowSize - curOnes

        for right in range(windowSize, n):
            curOnes += data[right]
            curOnes -= data[right - windowSize]
            res = min(res, windowSize - curOnes)

        return res


print(Solution().minSwaps([1, 0, 1, 0, 1]))
# 输出：1
# 解释：
# 有三种可能的方法可以把所有的 1 组合在一起：
# [1,1,1,0,0]，交换 1 次；
# [0,1,1,1,0]，交换 2 次；
# [0,0,1,1,1]，交换 1 次。
# 所以最少的交换次数为 1。
