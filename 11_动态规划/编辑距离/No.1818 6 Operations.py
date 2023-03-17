# https://yukicoder.me/problems/no/1818

# 6种操作
# 1. 将A[i]的值变为A[i]+1
# 2. 将A[i]的值变为A[i]-1
# 3. 将A[i]的值变为A[i]+A[i+1],删除A[i+1]
# 4. 将A[i]的值变为A[i]+A[i+1]+1,删除A[i+1]
# 5. 将A[i]的值变为A[i]-x,在A[i]后面插入x
# 6. 将A[i]的值变为A[i]-x-1,在A[i]后面插入x
# A种任何元素不能变为负数
# 求A变为B的最少操作次数
# !n,m<=1e3 sum(A)<=3e3 sum(B)<=3e3

# !将原数组的每个数变为[0,1,1,1,1...,1] (0+nums[i]个1)), 那么
# 1 => 插入(1)
# 2 => 删除(1)
# 3 => 删除(0)
# 4 => 替换(0->1)
# 5 => 替换(1->0)
# 6 => 替换(1->0) + 删除(1)

from 编辑距离 import editDistance

from typing import List


def operation6(nums1: List[int], nums2: List[int]) -> int:
    sb1, sb2 = [], []
    for v in nums1:
        sb1.append(0)
        for _ in range(v):
            sb1.append(1)
    for v in nums2:
        sb2.append(0)
        for _ in range(v):
            sb2.append(1)

    return editDistance(sb1, sb2)


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    print(operation6(nums1, nums2))
