# F - Erase between X and Y
# https://atcoder.jp/contests/abc421/tasks/abc421_f
# 第i个查询有两种类型：
# 1 x: 在 x 后面插入 i
# 2 x y: 删除 x 和 y 之间的元素，并输出这些元素的和
# https://atcoder.jp/contests/abc421/editorial/13787
#
# !关键：两个指针一起同时寻找

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    Q = int(input())

    next_node = [-1] * (Q + 1)

    for i in range(1, Q + 1):
        query = list(map(int, input().split()))
        q_type = query[0]

        if q_type == 1:
            x = query[1]
            next_node[i] = next_node[x]
            next_node[x] = i
        else:
            x, y = query[1], query[2]
            cx, cy = x, y
            sx, sy = 0, 0
            while True:
                if next_node[cx] != -1:
                    cx = next_node[cx]
                    if cx == y:
                        print(sx)
                        next_node[x] = y
                        break
                    sx += cx
                if next_node[cy] != -1:
                    cy = next_node[cy]
                    if cy == x:
                        print(sy)
                        next_node[y] = x
                        break
                    sy += cy
