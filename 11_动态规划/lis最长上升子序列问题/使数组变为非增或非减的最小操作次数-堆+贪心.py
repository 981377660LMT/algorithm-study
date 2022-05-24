# Make Array Non-decreasing or Non-increasing
# 每次操作可以使每个数加1或者减1
# 修改序列为非降或非增的最小修改次数

# 1 <= nums.length <= 1e5
# 0 <= nums[i] <= 1e9


from heapq import heappop, heappush, heappushpop, heapreplace
from typing import List

# https://codeforces.com/blog/entry/47821  slope trick
# https://www.luogu.com.cn/problem/solution/P4597

# https://github.dev/EndlessCheng/codeforces-go
# 修改序列为非降或非增的最小修改次数
# 单次修改可以把某个数 +1 或 -1
# https://www.luogu.com.cn/problem/solution/P4597
# 通过一个例子来解释这个基于堆的算法：1 5 10 4 2 2 2 2
# 假设当前维护的是非降序列，前三个数直接插入，不需要任何修改
# 插入 4 的时候，可以修改为 1 5 5 5，或 1 5 6 6，或... 1 5 10 10，修改次数均为 6
# 但我们也可以把修改后的序列视作 1 5 4 4，虽然序列不为非降序列，但修改的次数仍然为 6
# 接下来插入 2，基于 1 5 5 5 的话，修改后的序列就是 1 5 5 5 5，总的修改次数为 9
# 但我们也可以把修改后的序列视作 1 2 4 4 2，总的修改次数仍然为 9
# 接下来插入 2，如果基于 1 5 5 5 5 变成 1 5 5 5 5 5，会得到错误的修改次数 12
# 但是实际上有更优的修改 1 4 4 4 4 4，总的修改次数为 11
# 同上，把这个序列视作 1 2 2 4 2 2，总的修改次数仍然为 11
# !https://www.acwing.com/problem/content/description/275/


class Solution:
    def convertArray(self, nums: List[int]) -> int:
        def helper(nums: List[int]) -> int:
            """变为不减数组的最小操作次数
            
            如果num比前面的数小，那么就把前面的最大数变小
            """
            res, pq = 0, []  # 大根堆
            for num in nums:
                if not pq:
                    heappush(pq, -num)
                else:
                    preMax = -pq[0]
                    if preMax > num:
                        res += preMax - num
                        heappushpop(pq, -num)
                    heappush(pq, -num)
            return res

        return min(helper(nums), helper(nums[::-1]))


print(Solution().convertArray(nums=[3, 2, 4, 5, 0]))
print(Solution().convertArray(nums=[3, 1, 2, 1]))
print(Solution().convertArray([11, 11, 13, 8, 18, 19, 20, 7, 16, 3]))

