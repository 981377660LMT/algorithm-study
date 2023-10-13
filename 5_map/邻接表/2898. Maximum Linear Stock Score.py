# !求和最大的子序列，使得任意相邻两项之差等于对应的索引之差
# https://leetcode.cn/problems/maximum-linear-stock-score/


# arr[i]-arr[j] = i-j
# 变形为 arr[i]-i = arr[j]-j


from collections import defaultdict
from typing import List


class Solution:
    def maxScore(self, prices: List[int]) -> int:
        mp = defaultdict(int)
        for i, v in enumerate(prices):
            mp[v - i] += v
        return max(mp.values(), default=0)


if __name__ == "__main__":
    print(Solution().maxScore([1, 5, 3, 7, 8]))
