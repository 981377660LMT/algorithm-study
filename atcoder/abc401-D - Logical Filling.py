# D - Logical Filling
#
# https://atcoder.jp/contests/abc401/tasks/abc401_d
#
# - 给定包含 `.`、`o` 和 `?` 的字符串 S
# - 需要将所有 `?` 替换为 `.` 或 `o`，满足:
#   1. 恰好有 K 个 `o`
#   2. 没有连续的 `o`
# - 输出一个字符串 T，其中每个位置:
#   - 如果所有可能的有效替换该位置都是 `.`，则 T_i = `.`
#   - 如果所有可能的有效替换该位置都是 `o`，则 T_i = `o`
#   - 如果两种情况都存在，则 T_i = `?`


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, K = map(int, input().split())
    S = input()
