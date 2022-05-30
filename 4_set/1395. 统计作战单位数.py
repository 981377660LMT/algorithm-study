from typing import List

from sortedcontainers import SortedList

# 3 <= n <= 1000
# 作战单位需满足： rating[i] < rating[j] < rating[k] 或者 rating[i] > rating[j] > rating[k] ，其中  0 <= i < j < k < n

# nlogn


class Solution:
    def numTeams(self, rating: List[int]) -> int:
        left, right = SortedList(), SortedList(rating)
        res = 0
        for num in rating:
            right.discard(num)
            leftSmaller = left.bisect_left(num)
            rigthBigger = len(right) - right.bisect_left(num)
            res += leftSmaller * rigthBigger
            leftBigger = len(left) - left.bisect_left(num)
            rightSmaller = right.bisect_left(num)
            res += leftBigger * rightSmaller
            left.add(num)
        return res


print(Solution().numTeams(rating=[2, 5, 3, 4, 1]))
# 输出：3
# 解释：我们可以组建三个作战单位 (2,3,4)、(5,4,1)、(5,3,1) 。
