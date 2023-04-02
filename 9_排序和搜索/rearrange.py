# 将数组按照order里的顺序重新排列
from typing import List


def rearrage(A, order: List[int]):
    """将数组按照order里的顺序重新排序.
    A[order[0]], A[order[1]], ...
    """
    res = [None] * len(A)
    for i in range(len(order)):
        res[i] = A[order[i]]
    return res


# 6364. 老鼠和奶酪
# https://leetcode.cn/problems/mice-and-cheese/
class Solution:
    def miceAndCheese1(self, reward1: List[int], reward2: List[int], k: int) -> int:
        n = len(reward1)
        order = [i for i in range(n)]
        order.sort(key=lambda x: reward1[x] - reward2[x], reverse=True)
        res = 0
        for i in range(n):
            if i < k:
                res += reward1[order[i]]
            else:
                res += reward2[order[i]]
        return res

    def miceAndCheese2(self, reward1: List[int], reward2: List[int], k: int) -> int:
        select = [(a, b) for a, b in zip(reward1, reward2)]
        select.sort(key=lambda x: x[0] - x[1], reverse=True)
        return sum(a for a, _ in select[:k]) + sum(b for _, b in select[k:])
