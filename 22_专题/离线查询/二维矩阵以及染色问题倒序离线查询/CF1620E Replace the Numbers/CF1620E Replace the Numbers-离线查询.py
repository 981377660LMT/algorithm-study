# CF1620E Replace the Numbers-离线查询
# https://www.luogu.com.cn/problem/CF1620E
# 给出 q 个操作，操作分为两种：

# 1 x 在序列末尾插入数字 x。
# 2 x y 把序列中的所有 x 替换为 y。

# 求这个序列操作后的结果。

import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    q = int(input())
    queries = []
    nums = []
    for _ in range(q):
        t, *args = map(int, input().split())
        if t == 1:
            x = args[0]
            queries.append((1, len(nums), x))  # 索引赋值
            nums.append(x)
        else:
            x, y = args
            queries.append((2, x, y))

    last = dict()  # 终将成为你
    for i in range(q - 1, -1, -1):
        t, a, b = queries[i]
        if t == 1:
            nums[a] = last.get(b, b)
        else:
            last[a] = last.get(b, b)

    print(*nums)
