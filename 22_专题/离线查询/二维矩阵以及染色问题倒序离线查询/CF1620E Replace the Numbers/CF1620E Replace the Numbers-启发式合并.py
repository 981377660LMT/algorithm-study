# CF1620E Replace the Numbers-在线查询
# https://www.luogu.com.cn/problem/CF1620E
# 给出 q 个操作，操作分为两种：

# 1 x 在序列末尾插入数字 x。
# 2 x y 把序列中的所有 x 替换为 y。

# 求这个序列操作后的结果。

from collections import defaultdict
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    q = int(input())
    pos = defaultdict(list)  # !记录每个数字的下标
    count = 0

    def add(index: int, value: int) -> None:
        pos[value].append(index)

    def merge(from_: int, to: int) -> None:
        if from_ == to:
            return
        if len(pos[from_]) > len(pos[to]):
            pos[from_], pos[to] = pos[to], pos[from_]
        pos[to] += pos[from_]
        pos[from_] = []

    for _ in range(q):
        t, *args = map(int, input().split())
        if t == 1:
            x = args[0]
            add(count, x)
            count += 1
        else:
            x, y = args
            merge(x, y)

    res = [0] * count
    for k, v in pos.items():
        for i in v:
            res[i] = k
    print(*res)
