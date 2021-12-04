from typing import List
from bisect import bisect_left, bisect_right

INF = 2 ** 63 - 1

# 对于每个查询，你需要找到 子字符串中 在 两支蜡烛之间 的盘子的 数目 。
# '*' 表示一个 盘子 ，'|' 表示一支 蜡烛 。


class Solution:
    def platesBetweenCandles1(self, s: str, queries: List[List[int]]) -> List[int]:
        candles = [i for i, char in enumerate(s) if char == '|']
        res = []
        for left, right in queries:
            lower = bisect_left(candles, left)
            upper = bisect_right(candles, right) - 1
            res.append((candles[upper] - candles[lower]) - (upper - lower) if lower < upper else 0)
        return res

    # 预处理nextCandle 前缀处理
    def platesBetweenCandles(self, s: str, queries: List[List[int]]) -> List[int]:
        candelSum = [0] * (len(s) + 1)
        next = [INF] * (len(s) + 1)
        pre = [0] * (len(s) + 1)

        for i, char in enumerate(s):
            candelSum[i + 1] = candelSum[i] + int(char == '|')
            pre[i + 1] = i if char == '|' else pre[i]
        for i, char in reversed(list(enumerate(s))):
            next[i] = i if char == '|' else next[i + 1]

        res = []
        for left, right in queries:
            lower = next[left]
            upper = pre[right + 1]
            res.append(
                upper - lower - (candelSum[upper] - candelSum[lower]) if lower < upper else 0
            )

        return res


print(Solution().platesBetweenCandles(s="**|**|***|", queries=[[2, 5], [5, 9]]))
# 输出：[2,3]
# 解释：
# - queries[0] 有两个盘子在蜡烛之间。
# - queries[1] 有三个盘子在蜡烛之间。


# 方法一：二分

# 方法二：Next Candle O(n)
#                         0  1  2  3  4  5  6  7  8  9  10  11  12  13  14  15  16  17  18  19  20

#                         *  *  *  |  *  *  |  *  *  *   *   *   |   *   *   |   |   *   *   |   *
# nearest right candle:   3  3  3  3  6  6  6  12 12 12  12 12  12  15  15  15   16  19  19  19  -
# nearest left candle:    -  -  -  3  3  3  6  6  6  6   6  6   12  12  12  15  16  16  16   19  19
# candle count:           0  0  0  1  1  1  2  2  2  2   2  2    3   3   3   4   5   5   5   6   6
