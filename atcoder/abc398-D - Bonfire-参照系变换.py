# abc398-D - Bonfire - 参照系变换
# https://atcoder.jp/contests/abc398/tasks/abc398_d
#
# 题目描述
# 在一个无限扩展的二维网格上，网格的坐标 (0,0) 处有一个篝火。
# 在时刻 t = 0 时，只有 (0,0) 这个格子中存在烟雾。
#
# 给定一个由字符 N、W、S、E 组成的长度为 N 的字符串 S。
# 在时刻 t = 1, 2, …, N 依次发生以下现象：
#
# 1. 风吹来，使得当前存在的所有烟雾均会按照如下规则移动：
#   - 当 S 的第 t 个字符为 N 时，位于格子 (r, c) 的烟雾会移动到 (r − 1, c)；
#   - 当 S 的第 t 个字符为 W 时，位于格子 (r, c) 的烟雾会移动到 (r, c − 1)；
#   - 当 S 的第 t 个字符为 S 时，位于格子 (r, c) 的烟雾会移动到 (r + 1, c)；
#   - 当 S 的第 t 个字符为 E 时，位于格子 (r, c) 的烟雾会移动到 (r, c + 1)。
#
# 2. 如果格子 (0,0) 此时不含有烟雾，则在 (0,0) 生成新的烟雾。
#
# 高桥君站在格子 (R, C) 上。
#
# 对于每个整数 1 ≤ t ≤ N，判断在时刻 t + 0.5 时，格子 (R, C) 是否存在烟雾，并按照题目的要求输出结果。
#
# 约束条件
# - N 是 1 到 200,000 之间的整数
# - 字符串 S 由字符 N, W, S, E 组成，长度为 N
# - R, C 是介于 −N 和 N 之间的整数
# - (R, C) ≠ (0,0)
#
# !更换参考系 -> 以烟为参考系，人和篝火在移动，而不是烟在移动

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, R, C = map(int, input().split())
    S = input()

    dx, dy = 0, 0
    st = set({(0, 0)})
    res = []
    for c in S:
        if c == "N":
            dx -= 1
        if c == "W":
            dy -= 1
        if c == "S":
            dx += 1
        if c == "E":
            dy += 1
        st.add((-dx, -dy))
        pos = (R - dx, C - dy)
        res.append("1" if pos in st else "0")
    print("".join(res))
