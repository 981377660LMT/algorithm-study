# https://leetcode.cn/problems/maximum-number-of-intersections-on-the-chart/description/
# from collections import defaultdict

# class Solution:
#     def maxIntersectionCount(self, y: List[int]) -> int:
#         in_dict = defaultdict(int)
#         out_dict = defaultdict(int)
#         n = len(y)
#         cur_count = y.count(y[0])
#         for i in range(1, n):
#             min_y = min(y[i], y[i - 1])
#             max_y = max(y[i], y[i - 1])
#             in_dict[min_y] += 1
#             out_dict[max_y] += 1
#             if min_y < y[0] < max_y:
#                 cur_count += 1
#         lst = sorted(set(in_dict.keys()) | set(out_dict.keys()))
#         max_count = cur_count
#         cur_count = 0
#         for num in lst:
#             cur_count += in_dict[num] - out_dict[num]
#             max_count = max(max_count, cur_count)
#         return max_count

# 作者：Arnold
# 链接：https://leetcode.cn/problems/maximum-number-of-intersections-on-the-chart/solutions/2619799/python3sao-miao-xian-by-arnold-sb6ffylaw-wc1g/
# 来源：力扣（LeetCode）
# 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。


# https://leetcode.cn/problems/maximum-number-of-intersections-on-the-chart/description/
# 有一条由 n 个点连接而成的折线图。给定一个 下标从 1 开始 的整数数组 y，第 k 个点的坐标是 (k, y[k])。
# 图中没有水平线，即没有两个相邻的点有相同的 y 坐标。
# 假设在图中任意画一条无限长的水平线。请返回这条水平线与折线相交的最多交点数。
# 2 <= y.length <= 1e5
# 1 <= y[i] <= 1e9
# 对于范围 [1, n - 1] 内的所有 i，都有 y[i] != y[i + 1]


from typing import List, Tuple


def minmax(a: int, b: int) -> Tuple[int, int]:
    return (a, b) if a < b else (b, a)


class Solution:
    def maxIntersectionCount(self, ys: List[int]) -> int:
        n = len(ys)
