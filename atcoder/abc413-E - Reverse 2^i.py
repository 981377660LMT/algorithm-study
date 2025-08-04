# 给定长度为 2^N 的排列 P，
# 可以任意次对齐于 2^b 大小的区间执行区间反转操作，
# 求能够得到的字典序最小的排列。
#
# 把 P 看成完全二叉树的叶节点序列
# !问题变为：给定完全二叉树的叶标签（初始 P），对每个内部节点决定是否交换左右子叶序列，使得最终叶序列字典序最小。

import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e6))

input = lambda: sys.stdin.readline().rstrip("\r\n")


def solve(perm: List[int]) -> Tuple[int, ...]:
    def dfs(l: int, r: int) -> Tuple[int, ...]:
        if l + 1 == r:
            return (perm[l],)
        m = (l + r) // 2
        left, right = dfs(l, m), dfs(m, r)
        cand1, cand2 = left + right, right + left
        return min(cand1, cand2)

    return dfs(0, len(perm))


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        _ = int(input())
        P = list(map(int, input().split()))
        res = solve(P)
        print(" ".join(map(str, res)))
