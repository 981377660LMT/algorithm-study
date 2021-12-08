from typing import List

# 3 <= n <= 1000
# 作战单位需满足： rating[i] < rating[j] < rating[k] 或者 rating[i] > rating[j] > rating[k] ，其中  0 <= i < j < k < n

# 枚举中间点 O(n^2)，记录前面有多少个数比当前数小或大
# 可用树状数组优化到nlogn


class Solution:
    def numTeams(self, rating: List[int]) -> int:
        n = len(rating)
        biggerThanLeft = [0] * n
        smallerThanLeft = [0] * n

        res = 0
        for right in range(n):
            for mid in range(right):
                if rating[right] > rating[mid]:
                    biggerThanLeft[right] += 1
                    res += biggerThanLeft[mid]
                else:
                    smallerThanLeft[right] += 1
                    res += smallerThanLeft[mid]

        return res


print(Solution().numTeams(rating=[2, 5, 3, 4, 1]))
# 输出：3
# 解释：我们可以组建三个作战单位 (2,3,4)、(5,4,1)、(5,3,1) 。
