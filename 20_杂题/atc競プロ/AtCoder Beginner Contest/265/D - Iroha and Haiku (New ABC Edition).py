# 前缀和+转化式子
# 题目大意是给你一个下标从0开始长度为N的`正整数`数组A=(A0,A1,... .AN-1)
# 然后问你存不存在一个四元对(x,y,z,w)满足如下条件:
# 0<=x<y<z<w<=N
# Ax+...Ay-1=P
# Ay+...+Az-1=Q
# Az+...+Aw-1=R

# 固定x 就可以二分出y
# 之后就可以二分出z/w
# O(nlogn)
# !注意到数组是正整数数组 查找索引可以二分
# !用SortedList二分查找索引

from itertools import accumulate
from sortedcontainers import SortedList

N, P, Q, R = map(int, input().split())
nums = list(map(int, input().split()))
preSum = [0] + list(accumulate(nums))
sl = SortedList(preSum)

for i1 in range(N):
    try:
        i2 = sl.index(sl[i1] + P)
        i3 = sl.index(sl[i2] + Q)
        i4 = sl.index(sl[i3] + R)
        print("Yes")
        exit(0)
    except ValueError:
        continue

print("No")
