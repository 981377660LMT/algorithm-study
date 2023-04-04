# 给你一个整数 n 和一个在范围 [0, n - 1] 以内的整数 p ，
# 它们表示一个长度为 n 且下标从 0 开始的数组 arr ，
# 数组中除了下标为 p 处是 1 以外，其他所有数都是 0 。

# 同时给你一个整数数组 banned ，它包含数组中的一些位置。
# banned 中第 i 个位置表示 arr[banned[i]] = 0 ，题目保证 banned[i] != p 。

# 你可以对 arr 进行 若干次 操作。一次操作中，你选择大小为 k 的一个 子数组 ，
# 并将它 翻转 。在任何一次翻转操作后，你都需要确保 arr 中唯一的 1 不会到达任何 banned 中的位置。
# 换句话说，arr[banned[i]] 始终 保持 0 。

# 请你返回一个数组 ans ，对于 [0, n - 1] 之间的任意下标 i ，
# ans[i] 是将 1 放到位置 i 处的 最少 翻转操作次数，如果无法放到位置 i 处，此数为 -1 。

# !求从p位置出发到其他顶点的最短路
# !很多条边的最短路=>线段树优化建图/在线bfs(和求完全图的最小生成树算法类似)
# !这里线段树建图分奇偶比较麻烦,所以采用在线bfs
# python 在线bfs 解决边数很多的最短路问题
# https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/

from collections import deque
from typing import Callable, List, Optional, Tuple

INF = int(1e18)


class Solution:
    def minReverseOperations(self, n: int, p: int, banned: List[int], k: int) -> List[int]:
        def getRange(i: int) -> Tuple[int, int]:
            """反转长度为k的子数组,从i可以到达的左右边界闭区间."""
            return max(i - k + 1, k - i - 1), min(i + k - 1, 2 * n - k - i - 1)

        def setUsed(u: int) -> None:
            """将u位置标记为已经访问过."""
            finder[u & 1].erase(u)

        def findUnused(u: int) -> Optional[int]:
            """找到一个未使用的位置.如果不存在,返回None."""
            left, right = getRange(u)
            # pre = finder[(u + k + 1) & 1].prev(right)
            # if pre is not None and left <= pre <= right:
            #     return pre
            next_ = finder[(u + k + 1) & 1].next(left)
            if next_ is not None and left <= next_ <= right:
                return next_

        finder = [Finder(n) for _ in range(2)]
        for i in range(n):
            finder[i & 1].insert(i)
        for i in banned:
            finder[i & 1].erase(i)

        dist = onlineBfs(n, p, setUsed, findUnused)
        return [d if d != INF else -1 for d in dist]


def onlineBfs(
    n: int, start: int, setUsed: Callable[[int], None], findUnused: Callable[[int], Optional[int]]
) -> List[int]:
    """在线bfs。不预先给出图, 而是通过两个函数 setUsed 和 findUnused 来在线寻找边。

    Args:
        n (int): 顶点数。
        start (int): 起点。
        setUsed (Callable[[int], None]): 将 u 标记为已访问。
        findUnused (Callable[[int], int]): 找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `None`。

    Returns:
        List[int]: 从起点到各个点的距离。
    """
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    setUsed(start)
    while queue:
        cur = queue.popleft()
        while True:
            next_ = findUnused(cur)
            if next_ is None:
                break
            dist[next_] = dist[cur] + 1
            queue.append(next_)
            setUsed(next_)
    return dist


from typing import Optional


class Finder:
    """利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    __slots__ = "_n", "_lg", "_seg"

    @staticmethod
    def _trailingZeros1024(x: int) -> int:
        if x == 0:
            return 1024
        return (x & -x).bit_length() - 1

    def __init__(self, n: int) -> None:
        self._n = n
        seg = []
        while True:
            seg.append([0] * ((n + 1023) >> 10))
            n = (n + 1023) >> 10
            if n <= 1:
                break
        self._seg = seg
        self._lg = len(seg)

    def insert(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] |= 1 << (i & 1023)
            i >>= 10

    def erase(self, i: int) -> None:
        for h in range(self._lg):
            self._seg[h][i >> 10] &= ~(1 << (i & 1023))
            if self._seg[h][i >> 10]:
                break
            i >>= 10

    def next(self, i: int) -> Optional[int]:
        """返回x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            i = 0
        if i >= self._n:
            return
        seg = self._seg
        for h in range(self._lg):
            if i >> 10 == len(seg[h]):
                break
            d = seg[h][i >> 10] >> (i & 1023)
            if d == 0:
                i = (i >> 10) + 1
                continue
            i += self._trailingZeros1024(d)
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += self._trailingZeros1024(seg[g][i >> 10])
            return i

    def prev(self, i: int) -> Optional[int]:
        """返回x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        if i < 0:
            return
        if i >= self._n:
            i = self._n - 1
        seg = self._seg
        for h in range(self._lg):
            if i == -1:
                break
            d = seg[h][i >> 10] << (1023 - (i & 1023)) & ((1 << 1024) - 1)
            if d == 0:
                i = (i >> 10) - 1
                continue
            i += d.bit_length() - 1024
            for g in range(h - 1, -1, -1):
                i <<= 10
                i += (seg[g][i >> 10]).bit_length() - 1
            return i

    def islice(self, begin: int, end: int):
        """遍历[start,end)区间内的元素."""
        x = begin - 1
        while True:
            x = self.next(x + 1)
            if x is None or x >= end:
                break
            yield x

    def __contains__(self, i: int) -> bool:
        return self._seg[0][i >> 10] & (1 << (i & 1023)) != 0

    def __iter__(self):
        yield from self.islice(0, self._n)

    def __repr__(self):
        return f"FastSet({list(self)})"
