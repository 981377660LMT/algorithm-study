"""
简化版
1. 所有元素均为非负数时

2. 引入负数时
把负数变为绝对值 
此时只需要在答案上减去所有负数的绝对值即可
因为此时选负数等价于不选变化后的正数
不选负数对应选变化后的正数
"""

from typing import List
from sortedcontainers import SortedList

# 1 <= n <= 15
# 所有元素均为非负数：递推最小
# 所有元素可以为任意整数：平移-min(nums)


# 先给所有数加上 -min_element 转化成非负问题
# https://leetcode-cn.com/problems/find-array-given-subset-sums/solution/ti-jie-cong-zi-ji-de-he-huan-yuan-shu-zu-q9qw/
# https://leetcode-cn.com/problems/find-array-given-subset-sums/solution/jian-yi-ti-jie-by-sfiction-9i43/

# !从子集和还原数组
class Solution:
    def recoverArray(self, n: int, sums: List[int]) -> List[int]:
        OFFSET = -min(sums)
        sortedSums = SortedList()
        for s in sums:
            sortedSums.add(s + OFFSET)
        # 删去什么都不选的情况
        sortedSums.pop(0)

        res = []
        # 最小的数
        res.append(sortedSums[0])

        # 若我们已经推出了 res 中最小的 k 个元素，那么我们从 S 中把这 k 个元素所有子集的和去除
        for count in range(1, n):
            for state in range(1 << count):
                # 避免重复删除(保证res最新的一位被取到),即[0,1]删除了0+1那么[0,1,2]就不能删除0+1了
                if (state >> (count - 1)) & 1:
                    curSum = 0
                    for j in range(count):
                        if (state >> j) & 1:
                            curSum += res[j]
                    sortedSums.discard(curSum)
            res.append(sortedSums[0])

        # 处理负数
        # 如何找出原来的负数？=>原来的所有负数之和必定等于min(sums)，找出一个和为BIAS的子集，把他们转为负数即可
        for state in range(1 << n):
            curSum = 0
            for i in range(n):
                if (state >> i) & 1:
                    curSum += res[i]
            if curSum == OFFSET:
                for i in range(n):
                    if (state >> i) & 1:
                        res[i] *= -1
                break
        return res


print(Solution().recoverArray(n=3, sums=[-3, -2, -1, 0, 0, 1, 2, 3]))
print(Solution().recoverArray(n=2, sums=[0, 0, 0, 0]))
print(Solution().recoverArray(n=4, sums=[0, 0, 5, 5, 4, -1, 4, 9, 9, -1, 4, 3, 4, 8, 3, 8]))


# TODO
