# 给出n个巧克力和m个盒子的长宽，
# !每个盒子只能装一个巧克力而且不可旋转盒子，
# 问能否把巧克力都放进盒子里。

# !1889. 装包裹的最小浪费空间

from sortedcontainers import SortedList

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = list(map(int, input().split()))
    D = list(map(int, input().split()))
    chocolates = sorted(((a, b) for a, b in zip(A, B)), reverse=True)
    boxes = sorted(((c, d) for c, d in zip(C, D)), reverse=True)

    boxId = 0
    sl = SortedList()

    # !一个维度排序，维护另一个维度 二维偏序问题
    for cx, cy in chocolates:
        while boxId < len(boxes) and boxes[boxId][0] >= cx:
            sl.add(boxes[boxId][1])
            boxId += 1

        # 选一个y维度最接近的
        if not sl:
            print("No")
            return

        pos = sl.bisect_left(cy)
        if pos == len(sl):
            print("No")
            exit(0)
        sl.discard(sl[pos])

    print("Yes")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
