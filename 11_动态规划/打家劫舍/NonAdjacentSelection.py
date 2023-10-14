# https://atcoder.jp/contests/joisc2018/tasks/joisc2018_j


from typing import List
from itertools import accumulate
from heapq import heapify, heappop, heappush


class NonAdjacentSelection:
    """从数组不相邻选择k(0<=k<=(n+1//2))个数,最大化和/最小化和."""

    __slots__ = ("_n", "_nums", "_history", "_minimize", "_solved")

    def __init__(self, nums: List[int], minimize=False):
        self._n = len(nums)
        self._nums = nums
        self._history = []
        self._minimize = minimize
        self._solved = False

    def solve(self) -> List[int]:
        if self._minimize:
            self._nums = [-x for x in self._nums]

        nums, history = self._nums, self._history
        n = self._n
        rest = [True] * (n + 2)
        rest[0] = rest[n + 1] = False
        left = [i - 1 for i in range(n + 2)]
        right = [i + 1 for i in range(n + 2)]
        range_ = [(0, 0)] + [(i, i + 1) for i in range(n)] + [(0, 0)]
        val = [0] + nums + [0]
        pq = [(-val[i + 1], i + 1) for i in range(n)]
        heapify(pq)

        res = [0]
        while pq:
            add, i = heappop(pq)
            add = -add
            if not rest[i]:
                continue
            res.append(res[-1] + add)
            L = left[i]
            R = right[i]
            history.append(range_[i])
            if 1 <= L:
                right[left[L]] = i
                left[i] = left[L]
            if R <= n:
                left[right[R]] = i
                right[i] = right[R]
            if rest[L] and rest[R]:
                val[i] = val[L] + val[R] - val[i]
                heappush(pq, (-val[i], i))
                range_[i] = (range_[L][0], range_[R][1])
            else:
                rest[i] = False
            rest[L] = rest[R] = False

        if self._minimize:
            res = [-x for x in res]

        self._solved = True
        return res

    def restore(self, k: int) -> List[int]:
        """选择k个数,使得和最大/最小,返回选择的数的下标."""
        assert 0 <= k <= (self._n + 1) // 2, "k must be in [0,(n+1)//2]"
        if not self._solved:
            self.solve()
        diff = [0] * (self._n + 1)
        for a, b in self._history[:k]:
            diff[a] += 1
            diff[b] -= 1
            print(a, b, self._history)
        diff = list(accumulate(diff))
        return [i for i in range(self._n) if diff[i] & 1]


if __name__ == "__main__":
    nums = list(range(1, 11))
    NAS = NonAdjacentSelection(nums)
    print(NAS.solve())
    print(NAS.restore(2))

    # https://atcoder.jp/contests/joisc2018/tasks/joisc2018_j
    n = int(input())
    nums = [int(input()) for _ in range(n)]
    res = NonAdjacentSelection(nums).solve()
    print(*res[1:], sep="\n")
