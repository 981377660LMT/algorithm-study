# from titan_pylib.data_structures.segment_quadratic_division.segment_quadratic_division import SegmentQuadraticDivision
from typing import Union, Callable, TypeVar, Generic, Iterable
from functools import reduce
from itertools import chain

T = TypeVar("T")


class SegmentQuadraticDivision(Generic[T]):
    def __init__(self, n_or_a: Union[int, Iterable[T]], op: Callable[[T, T], T], e: T):
        if isinstance(n_or_a, int):
            self.n = n_or_a
            a = [e] * self.n
        else:
            a = list(n_or_a)
            self.n = len(a)
        self.op = op
        self.e = e
        self.size = int(self.n**0.5) + 1
        self.bucket_cnt = (self.n + self.size - 1) // self.size
        self.data = [a[k * self.size : (k + 1) * self.size] for k in range(self.bucket_cnt)]
        self.bucket_data = [reduce(self.op, v) for v in self.data]

    def prod(self, l: int, r: int) -> T:
        """Return op([l, r)). / 0 <= l <= r <= n / O(√N)"""
        assert 0 <= l <= r <= self.n
        if l == r:
            return self.e
        k1 = l // self.size
        k2 = r // self.size
        l -= k1 * self.size
        r -= k2 * self.size
        if k1 == k2:
            s = reduce(self.op, self.data[k1][l:r])
        else:
            s = self.e
            if l < len(self.data[k1]):
                s = reduce(self.op, self.data[k1][l:])
            if k1 + 1 < k2:
                s = (
                    reduce(self.op, self.bucket_data[k1 + 1 : k2])
                    if s == self.e
                    else reduce(self.op, self.bucket_data[k1 + 1 : k2], s)
                )
            if k2 < self.bucket_cnt and r > 0:
                s = (
                    reduce(self.op, self.data[k2][:r])
                    if s == self.e
                    else reduce(self.op, self.data[k2][:r], s)
                )
        return s

    def all_prod(self) -> T:
        """Return op([0, n)). / O(√N)"""
        return reduce(self.op, self.bucket_data)

    def __getitem__(self, indx: int) -> T:
        k = indx // self.size
        return self.data[k][indx - k * self.size]

    def __setitem__(self, indx: int, key: T):
        k = indx // self.size
        self.data[k][indx - k * self.size] = key
        self.bucket_data[k] = reduce(self.op, self.data[k])

    def __str__(self):
        return str(list(chain(*self.data)))

    def __repr__(self):
        return f"SegmentQuadraticDivision({self})"
