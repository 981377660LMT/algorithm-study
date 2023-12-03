# https://www.acwing.com/activity/content/code/content/1044156/
'''
问题做一下转换，把环拆开，组成长度为原来两倍的序列，把每个车站的油量和下一步的距离的差值
作为数组元素的数值，对这个长度是2n的数组求前缀和序列，其实题目要求的就是所有长度是n的滑动
窗口里面前缀和的最小值和窗口开始位置前面一个前缀和的差值的最小值是不是小于0，如果小于0，
说明中间n次移动至少有一次出现了没有油的情况，顺时针和逆时针都做一遍单调队列求滑动窗口最小值
的流程，两次结果有一次是成功的，就可以从该位置绕一圈回到原点
'''
# 安全：往后长为n的区间内所有前缀和>=0 即长度为n的区间的最小值>=0

# John 打算驾驶一辆汽车周游一个环形公路。
# 公路上总共有 n 个车站，每站都有若干升汽油（有的站可能油量为零），每升油可以让汽车行驶一千米。
# John 必须从某个车站出发，一直按顺时针（或逆时针）方向走遍所有的车站，并回到起点。
# 在一开始的时候，汽车内油量为零，John 每到一个车站就把该站所有的油都带上（起点站亦是如此），行驶过程中不能出现没有油的情况。
# 任务：判断以`每个车站为起点`能否按条件成功周游一周。

# 正向反向各一次操作
from collections import deque
from itertools import accumulate
from typing import List

n = int(input())
flag = [False] * (n + 1)
gas, cost = [0] * (n), [0] * (n)
for i in range(n):
    plus, minus = map(int, input().split())
    gas[i], cost[i] = plus, minus

# s1[1:n+1]:顺时针走每个加油站的油减去到下一站的油耗，那么s1的前缀和pre[i]就是到第i+1站剩下的油
# !这里取1～n为了配合计算前缀和
# pre[n]就是从站1回到站1所剩下的油，画个图就理解了。
nums1 = [gas[i] - cost[i] for i in range(n)]
# s2[1:n+1]:逆时针走每个加油站的油减去到下一站的油耗，那么s2的前缀和pre[i]就是到第i-1站剩下的油
# 用s2计算的时候反个向就刚好按照顺时针计算就可以了，返回结果的序号再处理一下就好了
nums2 = [gas[i] - cost[i - 1] for i in range(n)]


def cal(nums: List[int]) -> List[int]:
    res = [0] * (n + 1)
    nums = [*nums, *nums]
    preSum = [0] + list(accumulate(nums))
    queue = deque()
    # i不需要到最后一个元素，因为i==n的时候跟i==2×n的时候重复了。
    for i in range(1, 2 * n):
        # 队里始终保持n-1个元素
        if queue and i - queue[0] >= n:
            queue.popleft()
        if queue:
            if i >= n:
                # i从最后一站开始作为终点，如果这一圈里最小的剩油小于0，则i不能作为终点。
                if preSum[queue[0]] - preSum[i - n] < 0:
                    res[i - n + 1] = 0
                else:
                    res[i - n + 1] = 1
        while queue and preSum[queue[-1]] >= preSum[i]:
            queue.pop()
        queue.append(i)

    return res[1:]


res1 = cal(nums1)
res2 = cal(nums2[::-1])[::-1]


for r1, r2 in zip(res1, res2):
    if r1 or r2:
        print('TAK')
    else:
        print('NIE')


# TODO
