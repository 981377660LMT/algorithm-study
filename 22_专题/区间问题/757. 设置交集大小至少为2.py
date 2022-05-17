from typing import List

# 一个整数区间 [a, b]  ( a < b ) 代表着从 a 到 b 的所有连续整数，包括 a 和 b。
# 给你一组整数区间intervals，请找到一个最小的集合 S，
# 使得 S 里的元素与区间intervals中的每一个整数区间都至少有2个元素相交。
# 输出这个最小集合S的大小。


# 总结:区间相交问题：排序
# 优先处理长度小的区间，比较curStart与preEnd
class Solution:
    def intersectionSizeTwo(self, intervals: List[List[int]]) -> int:
        # 当前窗口
        res, left, right = 0, -1, -1

        # end相同时，越往后走范围就向前覆盖得越大，即更加符合香蕉条件
        intervals = sorted(intervals, key=lambda x: (x[1], -x[0]))

        for s, e in intervals:
            if s > right:  # <..end1,..end2,..s..e>
                left, right = e - 1, e
                res += 2
            elif s > left:  # <..end1,.s..end2..e>
                left, right = right, e
                res += 1

        return res


print(Solution().intersectionSizeTwo(intervals=[[1, 3], [1, 4], [2, 5], [3, 5]]))
# 输出: 3
# 解释:
# 考虑集合 S = {2, 3, 4}. S与intervals中的四个区间都有至少2个相交的元素。
# 且这是S最小的情况，故我们输出3。
