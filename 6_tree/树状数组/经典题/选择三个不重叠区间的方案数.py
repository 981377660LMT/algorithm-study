# three non-overlapping intervals
# 三个不重叠区间的方案数
# n<=1e5
# 1<=starting[i],ending[i]<=1e9

from typing import List

from sortedcontainers import SortedList

# 先按 starting 排序。

# 再对每一个区间，求左边不与该区间重叠的区间个数
# (即左边 ending[l] < starting[i] 的个数)，
# 和右边不重叠的区间个数，乘起来，并累加到答案中。


def getThreeNonOverlappingIntervals(intervals: List[List[int]]) -> int:
    """
    从n个区间(闭区间)中选出3个不重叠的区间,计算有多少种选择方式,返回数量.
    区间都是闭区间

    !枚举每个区间作为中间的区间
    """
    intervals.sort()
    leftEnds, rigthStarts = SortedList(), SortedList([s for s, _ in intervals])
    res = 0
    for start, end in intervals:
        rigthStarts.remove(start)
        leftSmaller = leftEnds.bisect_left(start)
        rightBigger = len(rigthStarts) - rigthStarts.bisect_right(end)
        res += leftSmaller * rightBigger
        leftEnds.add(end)
    return res


# starting=[1,2,4,3,7]
# ending=[3,4,6,5,8]
assert getThreeNonOverlappingIntervals([[1, 3], [2, 4], [3, 5], [4, 6], [7, 8]]) == 1


# starting=[5,2,3,7]
# ending=[7,2,4,8]
assert getThreeNonOverlappingIntervals([[5, 7], [2, 2], [3, 4], [7, 8]]) == 2
