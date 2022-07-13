# https://atcoder.jp/contests/abc253/submissions/33192417
# !insert 和 delete 操作的时间复杂度是O(logN)


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
    def __init__(self, iterable=None):
        self.root = None
        if iterable is not None:
            for i in iterable:
                self.insert(len(self), i)

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

    def __repr__(self) -> str:
        return str(self.root.val) if self.root is not None else "None"


if __name__ == "__main__":
    nums = SplayTreeList(range(100))
    nums.delete(5)
    print(len(nums))
    print(nums[5])
    nums.insert(5, 500)
    print(nums[5])
    print(nums[6])
