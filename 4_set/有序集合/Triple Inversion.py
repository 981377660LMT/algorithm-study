#  return the number of pairs i < j such that nums[i] > nums[j] * 3.
from sortedcontainers import SortedList


class Solution:
    def solve(self, A):
        sortedList = SortedList()
        res = 0

        for num in A:
            i = sortedList.bisect_right(num * 3)
            res += len(sortedList) - i
            sortedList.add(num)

        return res
