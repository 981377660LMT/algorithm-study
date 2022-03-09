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

n = int(input())
nums = [0] * n
for i in range(n):
    plus, minus = map(int, input().split())
    nums[i] = plus - minus


def check() -> None:
    queue = deque()
    for i in range(len(preSum)):
        while queue and i - queue[0] > n:
            queue.popleft()
        if i >= n - 1 and queue and preSum[queue[0]] >= 0:
            res[i % n] = True
        while queue and preSum[i] <= preSum[queue[-1]]:
            queue.pop()
        queue.append(i)


res = [False] * n

preSum = list(accumulate(nums * 2))
check()
preSum = preSum[::-1]
check()

for i in range(n):
    if res[i]:
        print('TAK')
    else:
        print('NIE')


# 晕 太复杂了
