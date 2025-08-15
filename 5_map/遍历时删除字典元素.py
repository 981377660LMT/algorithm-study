from collections import defaultdict
from heapq import heapify, heappop
import random
from sortedcontainers import SortedDict

# !1.堆的删除与遍历通过while循环实现
pq = [random.randrange(1, 100) for _ in range(10)]
heapify(pq)

# !删除堆中小于10的元素
while pq and pq[0] < 10:
    heappop(pq)

###########################################################
# !2.遍历字典删除key
mp = {1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 0: 0, -1: -1}

# !list(mp)获取键的拷贝,注意不要边遍历边修改字典
for key in list(mp):
    if key <= 0:
        print(mp.pop(key))

###########################################################
# !3.删除SortedDict中小于等于0的key
sd = SortedDict({2: 2, 1: 1, 3: 3, 4: 4, 5: 5, 0: 0, -1: -1})
while sd and sd.peekitem(0)[0] <= 0:
    print(sd.popitem(0))
