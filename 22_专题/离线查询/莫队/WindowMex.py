# https://zhuanlan.zhihu.com/p/563177000 求区间mex的多种方法

# E - Mex Min-滑动窗口mex
# !O(nlogn) 定长滑动窗口mex (区间mex特殊情况)
# !（1）添加一个数到集合中
# !（2）从集合中删除一个数 (如果没有删除操作,只用while维护mex)
# !（3）查询这个集合的mex
# 1<=m<=n<=2e6
# 0<=nums[i]<n


from typing import List
from sortedcontainers import SortedList


class WindowMex:
    __slots__ = "_maxOperation", "_mexStart", "_counter", "_sl"

    def __init__(self, maxOperation: int, mexStart=0) -> None:
        self._maxOperation = maxOperation
        self._mexStart = mexStart
        self._counter = [0] * (maxOperation + 1)
        self._sl = SortedList(list(range(mexStart, mexStart + maxOperation + 1)))

    def add(self, v: int) -> bool:
        mexStart, maxOperation = self._mexStart, self._maxOperation
        if v < mexStart or v > mexStart + maxOperation:
            return False
        self._counter[v - mexStart] += 1
        if self._counter[v - mexStart] == 1:
            self._sl.remove(v)
        return True

    def discard(self, v: int) -> bool:
        mexStart, maxOperation = self._mexStart, self._maxOperation
        if v < mexStart or v > mexStart + maxOperation:
            return False
        if self._counter[v - mexStart] == 0:
            return False
        self._counter[v - mexStart] -= 1
        if self._counter[v - mexStart] == 0:
            self._sl.add(v)
        return True

    def query(self) -> int:
        if not self._sl:
            return self._mexStart
        return self._sl[0]  # type: ignore


def windowMex(nums: List[int], k: int) -> List[int]:
    """记长为k的滑动窗口的mex为mexi,求n-k+1个mexi的最小值"""

    def add(num: int) -> None:
        counter[num] += 1
        if counter[num] == 1:
            sl.remove(num)

    def remove(num: int) -> None:
        counter[num] -= 1
        if counter[num] == 0:
            sl.add(num)

    def query() -> int:
        return sl[0]  # type: ignore

    n = len(nums)
    res = []
    counter = [0] * (n + 10)
    sl = SortedList(list(range(n + 1)))  # 维护mex候选人0-n
    for right in range(n):
        add(nums[right])
        if right >= k:
            remove(nums[right - k])
        if right >= k - 1:
            res.append(query())
    return res


n, k = map(int, input().split())
nums = list(map(int, input().split()))
print(min(windowMex(nums, k)))
