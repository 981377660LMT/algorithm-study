from typing import List


class BinaryTrie:
    """
    Reference:
     - https://atcoder.jp/contests/arc028/submissions/19916627
     - https://judge.yosupo.jp/submission/35057
    """

    __slots__ = (
        "max_log",
        "x_end",
        "v_list",
        "multiset",
        "add_query_count",
        "add_query_limit",
        "edges",
        "size",
        "is_end",
        "max_v",
        "lazy",
    )

    def __init__(
        self,
        max_log=30,
        allow_multiple_elements=True,
        add_query_limit=10**6,
    ):
        self.max_log = max_log
        self.x_end = 1 << max_log
        self.v_list = [0] * (max_log + 1)
        self.multiset = allow_multiple_elements
        self.add_query_count = 0
        self.add_query_limit = add_query_limit
        n = max_log * add_query_limit + 1
        self.edges = [-1] * (2 * n)
        self.size = [0] * n
        self.is_end = [0] * n
        self.max_v = 0
        self.lazy = 0

    def add(self, x: int):
        # assert 0 <= x < self.x_end
        # assert 0 <= self.add_query_count < self.add_query_limit
        x ^= self.lazy
        v = 0
        for i in range(self.max_log - 1, -1, -1):
            d = (x >> i) % 2
            if self.edges[2 * v + d] == -1:
                self.max_v += 1
                self.edges[2 * v + d] = self.max_v
            v = self.edges[2 * v + d]
            self.v_list[i] = v
        if self.multiset or self.is_end[v] == 0:
            self.is_end[v] += 1
            for v in self.v_list:
                self.size[v] += 1
        self.add_query_count += 1

    def discard(self, x: int):
        if not 0 <= x < self.x_end:
            return
        x ^= self.lazy
        v = 0
        for i in range(self.max_log - 1, -1, -1):
            d = (x >> i) % 2
            if self.edges[2 * v + d] == -1:
                return
            v = self.edges[2 * v + d]
            self.v_list[i] = v
        if self.is_end[v] > 0:
            self.is_end[v] -= 1
            for v in self.v_list:
                self.size[v] -= 1

    def erase(self, x: int, count: int = -1):
        """删除count个x x=-1表示删除所有x"""
        # assert -1 <= count
        if not 0 <= x < self.x_end:
            return
        x ^= self.lazy
        v = 0
        for i in range(self.max_log - 1, -1, -1):
            d = (x >> i) % 2
            if self.edges[2 * v + d] == -1:
                return
            v = self.edges[2 * v + d]
            self.v_list[i] = v
        if count == -1 or self.is_end[v] < count:
            count = self.is_end[v]
        if self.is_end[v] > 0:
            self.is_end[v] -= count
            for v in self.v_list:
                self.size[v] -= count

    def count(self, x: int) -> int:
        if not 0 <= x < self.x_end:
            return 0
        x ^= self.lazy
        v = 0
        for i in range(self.max_log - 1, -1, -1):
            d = (x >> i) % 2
            if self.edges[2 * v + d] == -1:
                return 0
            v = self.edges[2 * v + d]
        return self.is_end[v]

    def bisect_left(self, x: int) -> int:
        if x < 0:
            return 0
        if self.x_end <= x:
            return len(self)
        v = 0
        ret = 0
        for i in range(self.max_log - 1, -1, -1):
            d = (x >> i) % 2
            l = (self.lazy >> i) % 2
            lc = self.edges[2 * v]
            rc = self.edges[2 * v + 1]
            if l == 1:
                lc, rc = rc, lc
            if d:
                if lc != -1:
                    ret += self.size[lc]
                if rc == -1:
                    return ret
                v = rc
            else:
                if lc == -1:
                    return ret
                v = lc
        return ret

    def bisect_right(self, x: int) -> int:
        return self.bisect_left(x + 1)

    def index(self, x: int) -> int:
        if x not in self:
            raise ValueError(f"{x} is not in BinaryTrie")
        return self.bisect_left(x)

    def find(self, x: int) -> int:
        if x not in self:
            return -1
        return self.bisect_left(x)

    def kth_elem(self, k: int) -> int:
        if k < 0:
            k += self.size[0]
        # assert 0 <= k < self.size[0]
        v = 0
        ret = 0
        for i in range(self.max_log - 1, -1, -1):
            l = (self.lazy >> i) % 2
            lc = self.edges[2 * v]
            rc = self.edges[2 * v + 1]
            if l == 1:
                lc, rc = rc, lc
            if lc == -1:
                v = rc
                ret |= 1 << i
                continue
            if self.size[lc] <= k:
                k -= self.size[lc]
                v = rc
                ret |= 1 << i
            else:
                v = lc
        return ret

    def minimum(self) -> int:
        return self.kth_elem(0)

    def maximum(self) -> int:
        return self.kth_elem(-1)

    def xor_all(self, x: int):
        # assert 0 <= x < self.x_end
        self.lazy ^= x

    def __iter__(self):
        q = [(0, 0)]
        for i in range(self.max_log - 1, -1, -1):
            l = (self.lazy >> i) % 2
            nq = []
            for v, x in q:
                lc = self.edges[2 * v]
                rc = self.edges[2 * v + 1]
                if l == 1:
                    lc, rc = rc, lc
                if lc != -1:
                    nq.append((lc, 2 * x))
                if rc != -1:
                    nq.append((rc, 2 * x + 1))
            q = nq
        for v, x in q:
            for _ in range(self.is_end[v]):
                yield x

    def __str__(self):
        prefix = "BinaryTrie("
        content = list(map(str, self))
        suffix = ")"
        if content:
            content[0] = prefix + content[0]
            content[-1] = content[-1] + suffix
        else:
            content = [prefix + suffix]
        return ", ".join(content)

    def __getitem__(self, k):
        return self.kth_elem(k)

    def __contains__(self, x: int) -> bool:
        return not not self.count(x)

    def __len__(self):
        return self.size[0]

    def __bool__(self):
        return not not len(self)

    def __ixor__(self, x: int):
        self.xor_all(x)
        return self


if __name__ == "__main__":
    # import sys

    # input = sys.stdin.readline

    # q = int(input())
    # bt = BinaryTrie(30, False, q)
    # res = []
    # for _ in range(q):
    #     t, x = map(int, input().split())
    #     if t == 0:
    #         bt.add(x)
    #     if t == 1:
    #         bt.discard(x)
    #     if t == 2:
    #         bt.xor_all(x)
    #         res.append(bt.minimum())  # 求x与树中异或最小值
    #         bt.xor_all(x)

    # print("\n".join(map(str, res)))

    bt = BinaryTrie(30, True, 10)
    bt.add(20)
    bt.add(1)
    bt.add(2)
    bt.add(21)
    print(bt, bt.bisect_left(2), bt.find(29))
    bt.add(2)
    bt.add(2)
    bt.erase(2, -1)
    print(bt, bt.bisect_left(2), bt.find(29))

    # 1803. 统计异或值在范围内的数对有多少
    class Solution:
        def countPairs(self, nums: List[int], low: int, high: int) -> int:
            n = len(nums)
            max_log = max(nums).bit_length()
            bt = BinaryTrie(add_query_limit=n, max_log=max_log, allow_multiple_elements=True)
            for num in nums:
                bt.add(num)
            res = 0
            for num in nums:
                bt.xor_all(num)
                res += bt.bisect_right(high) - bt.bisect_left(low)
                bt.xor_all(num)
            return res // 2
