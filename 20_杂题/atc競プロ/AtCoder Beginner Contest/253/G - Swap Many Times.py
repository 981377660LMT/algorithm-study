"""伸展树"""

import sys
import os


class STNode:
    __slots__ = "val", "l", "r", "p", "size"

    def __init__(self, val):
        self.val = val
        self.l, self.r, self.p = None, None, None
        self.size = 1

    def update(self):
        self.size = 1
        if self.l is not None:
            self.size += self.l.size
        if self.r is not None:
            self.size += self.r.size

    def state(self):
        if self.p is None:
            return 0
        if self.p.l == self:
            return 1
        else:
            return -1  # self.p.r == self

    def splay(self):
        while self.p is not None:
            st, st_p = self.state(), self.p.state()
            if st_p == 0:
                if st == 1:
                    self.p._rotr()
                else:
                    self.p._rotl()
            elif st_p == 1:
                if st == 1:
                    self.p.p._rotr()
                    self.p._rotr()
                else:
                    self.p._rotl()
                    self.p._rotr()
            else:
                if st == 1:
                    self.p._rotr()
                    self.p._rotl()
                else:
                    self.p.p._rotl()
                    self.p._rotl()
        return self

    def _rotl(self):
        nd = self.r
        nd.p = self.p
        if nd.p is not None:
            if nd.p.l is self:
                nd.p.l = nd
            else:
                nd.p.r = nd
        self.r = nd.l
        if self.r is not None:
            self.r.p = self
        self.p = nd
        nd.l = self
        self.update()
        nd.update()

    def _rotr(self):
        nd = self.l
        nd.p = self.p
        if nd.p is not None:
            if nd.p.r is self:
                nd.p.r = nd
            else:
                nd.p.l = nd
        self.l = nd.r
        if self.l is not None:
            self.l.p = self
        self.p = nd
        nd.r = self
        self.update()
        nd.update()


class SplayTreeList:
    def __init__(self):
        self.root = None

    def insert(self, index: int, val: int) -> None:
        rtl, rtr = self._split(index)
        rt = STNode(val)
        self.root = self._merge(self._merge(rtl, rt), rtr)

    def delete(self, index: int) -> None:
        self._splay(index)
        rtl = self.root.l
        rtr = self.root.r
        if rtl is not None:
            rtl.p = None
        if rtr is not None:
            rtr.p = None
        self.root = self._merge(rtl, rtr)

    def _splay(self, idx):
        nd = self.root
        while True:
            sizel = nd.l.size if nd.l is not None else 0
            if idx < sizel:
                nd = nd.l
            elif idx > sizel:
                nd = nd.r
                idx -= sizel + 1
            else:
                self.root = nd.splay()
                break

    def _merge(self, rtl, rtr):
        if rtl is None:
            return rtr
        if rtr is None:
            return rtl
        while rtl.r is not None:
            rtl = rtl.r
        rtl = rtl.splay()
        rtl.r = rtr
        rtr.p = rtl
        rtl.update()
        return rtl

    def _split(self, cntl):
        if cntl == 0:
            return None, self.root
        if cntl == self.root.size:
            return self.root, None
        self._splay(cntl)
        rtl = self.root.l
        rtr = self.root
        rtl.p = None
        rtr.l = None
        rtr.update()
        return rtl, rtr

    def __getitem__(self, idx):
        self._splay(idx)
        return self.root.val

    def __setitem__(self, idx, val):
        self._splay(idx)
        self.root.val = val

    def __len__(self):
        return self.root.size if self.root is not None else 0


# 对于n ，初始化一个序列a1,…, an ，满足a = i。(n<=2e5)
# 对于n，有n(n+1)个形如(x,y)的满足1≤a <y≤n 的数对，按照pair的规则排序。 ' 2
# 给定L,R，对于这个pair序列的第工个到第R个，依次操作:交换az和agy求最终的序列。

# !枚举x，然后只需要一个支持某个位置插入删除的数据结构就可以了。 (splay)
# https://www.zhihu.com/search?type=content&q=atcoder%20253
# https://atcoder.jp/contests/abc253/submissions/33192417


# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# MOD = int(1e9 + 7)


# def main() -> None:

#     n, left, right = map(int, input().split())
#     nums = [i + 1 for i in range(n)]

#     stl = SplayTreeList()
#     for i, v in enumerate(nums):
#         stl.insert(i, v)

#     times = 0
#     for pivot in range(n):
#         if left <= times + 1 and times + (n - pivot - 1) <= right:
#             stl.insert(pivot, stl[n - 1])
#             stl.delete(n)
#             times += n - pivot - 1
#         elif times + (n - pivot - 1) < left or right < times + 1:
#             times += n - pivot - 1
#             continue
#         else:
#             for i in range(pivot + 1, n):
#                 if left <= times + 1 <= right:
#                     stl[pivot], stl[i] = stl[i], stl[pivot]
#                 times += 1
#     print(*[stl[i] for i in range(n)])


# if __name__ == "__main__":
#     if os.environ.get("USERNAME", " ") == "caomeinaixi":
#         while True:
#             main()
#     else:
#         main()
