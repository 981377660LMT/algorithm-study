# No.2360 Path to Integer - 树上路径分数之和
# https://yukicoder.me/problems/no/2360
#
# 给定一棵树，树上每个节点都有一个数字，每一位都是0-9之间的整数，长度不超过9.
# 路径的分数定义为路径上所有节点的数字拼接成的整数.
# 请你求出树上所有路径的分数之和模998244353.

from Rerooting import Rerooting

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


from typing import Tuple

INF = int(4e18)
MOD = 998244353


pow10 = [1]
for _ in range(100):
    pow10.append(pow10[-1] * 10 % MOD)


if __name__ == "__main__":
    n = int(input())
    nums = [int(x) for x in input().split()]
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    pows = [pow10[len(str(nums[i]))] for i in range(n)]

    E = Tuple[int, int]  # (count,sum) 节点个数，路径和

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        count = childRes1[0] + childRes2[0]
        sum_ = childRes1[1] + childRes2[1]
        if sum_ >= MOD:
            sum_ -= MOD
        return (count, sum_)

    def composition(fromRes: E, parent: int, child: int, direction: int) -> E:
        """direction: 0: child -> parent, 1: parent -> child"""
        count, sum_ = fromRes
        from_ = child if direction == 0 else parent
        value, pow_ = nums[from_], pows[from_]
        return (count + 1, (sum_ * pow_ + value * (count + 1)) % MOD)

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e, op, composition)  # !注意dp不包含根节点
    res = 0
    for i in range(n):
        c, s = dp[i]
        v, p = nums[i], pows[i]
        res += s * p + v * (c + 1)
        res %= MOD
    print(res % MOD)
