# 我们有一个空序列A。请依次处理Q个命令，每个命令有三种类型，每种类型的格式如下:
# 1 x:将x加入A(不去重)
# 2 x k:求在A的≤x的元素中，第k大的值
# 3 x k:求在A的≥ x的元素中，第k小的值。

import sys
import os

from sortedcontainers import SortedList

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    q = int(input())
    sl = SortedList()
    for _ in range(q):
        t, *rest = map(int, input().split())
        if t == 1:
            x = rest[0]
            sl.add(x)
        elif t == 2:
            x, k = rest
            pos = sl.bisect_right(x) - 1
            kth = pos - (k - 1)
            if kth < 0:
                print(-1)
            else:
                print(sl[kth])
        elif t == 3:
            x, k = rest
            pos = sl.bisect_left(x)
            kth = pos + (k - 1)
            if kth >= len(sl):
                print(-1)
            else:
                print(sl[kth])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
