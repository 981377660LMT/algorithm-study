#  6340. 统计上升四元组
#  !i1 i2 i3 i4 满足的数需要为 1 3 2 4 大小关系(1324模式)
#  求满足条件的四元组的个数 n<=4000

#  !统计三元组:枚举中间的数,然后统计左侧比它小的数的个数,右侧比它大的数的个数
#  !统计四元组:枚举中间的两个数,然后统计左侧比它小的数的个数,右侧比它大的数的个数
#  类似2242. 节点序列的最大得分

#  !预处理:O(n^2)
#  !总结:枚举中间的数
#  1. 预处理 O(n^2)
#  2. 树状数组/SortedList O(n^2logn) TLE


from typing import List
from sortedcontainers import SortedList


class Solution:
    def countQuadruplets(self, nums: List[int]) -> int:
        n = len(nums)
        leftSmaller = SortedList()
        res = 0

        for i2 in range(n):
            num2 = nums[i2]
            rightBigger = SortedList(nums[i2 + 1 :])

            for i3 in range(i2 + 1, n):
                num3 = nums[i3]
                rightBigger.remove(num3)
                if num2 <= num3:
                    continue
                count1 = leftSmaller.bisect_left(num3)  # 统计i2左侧严格小于num3的数的个数
                count2 = len(rightBigger) - rightBigger.bisect_right(num2)  # 统计i2右侧严格大于num2的数的个数
                res += count1 * count2

            leftSmaller.add(num2)

        return res


print(Solution().countQuadruplets([1, 3, 2, 4, 5]))
