# 从一个单词中取一个字母所需要的代币数量，为该字母左边和右边字母数量之积

from functools import lru_cache
from collections import deque
from typing import Tuple


INF = int(1e12)


@lru_cache(None)
def calCost(need: Tuple[bool, ...]) -> int:
    """计算从两端取出所需字符的代价 deque"""
    n = len(need)
    queue = deque([i for i in range(n) if need[i]])
    cost = 0
    leftMoved, rightMoved = 0, 0
    while queue:
        left, right = queue[0], queue[-1]
        leftCost = (left - leftMoved) * (n - 1 - left - rightMoved)
        rightCost = (right - leftMoved) * (n - 1 - right - rightMoved)
        if leftCost <= rightCost:
            queue.popleft()
            cost += leftCost
            leftMoved += 1
        else:
            queue.pop()
            cost += rightCost
            rightMoved += 1
    return cost


@lru_cache(None)
def calCost2(need: Tuple[bool, ...]) -> int:
    """计算从两端取出所需字符的代价 双指针"""
    n = len(need)
    left, right = 0, n - 1
    cost = 0
    leftMoved, rightMoved = 0, 0
    while left <= right:
        while left <= right and not need[left]:
            left += 1
        while left <= right and not need[right]:
            right -= 1
        if left > right:
            break
        leftCost = (left - leftMoved) * (n - 1 - left - rightMoved)
        rightCost = (right - leftMoved) * (n - 1 - right - rightMoved)
        if leftCost <= rightCost:
            left += 1
            cost += leftCost
            leftMoved += 1
        else:
            right -= 1
            cost += rightCost
            rightMoved += 1
    return cost


print(calCost(tuple([True, False, True, False, True])))
print(calCost2(tuple([True, False, True, False, True])))
