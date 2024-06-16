# F - MST Query
# https://atcoder.jp/contests/abc355/tasks/abc355_f
import sys
from typing import List


input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    N, L, R = map(int, input().split())
    M = 1 << N
    offset = 1
    while offset < M:
        offset <<= 1

    def divide(start: int, end: int) -> List[int]:
        seg = []
        start += offset
        end += offset
        while start < end:
            if start & 1:
                seg.append(start)
                start += 1
            if end & 1:
                end -= 1
                seg.append(end)
            start >>= 1
            end >>= 1
        return seg

    def query(a: int, b: int) -> int:
        print(f"? {a} {b}", flush=True)
        return int(input())

    def output(res: int) -> None:
        print(f"! {res}", flush=True)

    def leftLeaf(x: int) -> int:
        while x < offset:
            x <<= 1
        return x - offset

    def rightLeaf(x: int) -> int:
        while x < offset:
            x = (x << 1) | 1
        return x - offset

    res = []
    for v in divide(L, R + 1):
        left, right = leftLeaf(v), rightLeaf(v)
        i = (right - left + 1).bit_length() - 1
        j = left // 2**i
        res.append(query(i, j))
    output(sum(res) % 100)
