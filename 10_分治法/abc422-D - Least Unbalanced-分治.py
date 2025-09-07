# abc422-D - Least Unbalanced-分治
# https://atcoder.jp/contests/abc422/tasks/abc422_d
# 设 N 为正整数。对于一个长度为 2^N 的非负整数序列 A = (A_1, A_2, ..., A_{2^N})，其不平衡度定义为通过以下操作得到的非负整数值：

# 首先，设 X = 0。
# 执行以下一系列操作 N 次： a. 更新 X：X = max(X, max(A) - min(A))。这里的 max(A) 和 min(A) 分别表示当前序列 A 的最大值和最小值。 b. 将 A 中从头开始每两个元素分为一组，将它们的和组成一个新序列，其长度是原序列 A 的一半。这个新序列成为新的 A。即 A ← (A_1 + A_2, A_3 + A_4, ..., A_{|A|-1} + A_{|A|})。
# 最终的 X 值就是不平衡度。
# 例如，对于 N=2, A=(6, 8, 3, 5)，不平衡度为 6，计算过程如下：

# 初始时, X=0。
# 第 1 轮操作:
# X 更新为 max(X, max(A) - min(A)) = max(0, 8 - 3) = 5。
# A 变为 (6+8, 3+5) = (14, 8)。
# 第 2 轮操作:
# X 更新为 max(X, max(A) - min(A)) = max(5, 14 - 8) = 6。
# A 变为 (14+8) = (22)。
# 最终 X 变为 6。
# 现在给定非负整数 K。请你构造一个长度为 2^N、所有元素总和为 K 的非负整数序列，使其不平衡度最小。

# 限制条件

# 1 ≤ N ≤ 20
# 0 ≤ K ≤ 10^9
# N, K 是整数
#
# !每一轮操作中的序列元素尽可能地接近。
# https://atcoder.jp/contests/abc422/editorial/13840

from typing import List


def dfs(start: int, end: int, div: int, mod: int, res: List[int]) -> None:
    if start >= end:
        return
    if start + 1 == end:
        res[start] = div + (1 if mod > 0 else 0)
        return
    mid = (start + end) // 2
    leftMod = mod // 2
    rightMod = mod - leftMod
    dfs(start, mid, div, leftMod, res)
    dfs(mid, end, div, rightMod, res)


if __name__ == "__main__":
    N, K = map(int, input().split())

    len_ = 1 << N
    div, mod = divmod(K, len_)
    unbalance = 1 if mod else 0
    print(unbalance)

    res = [0] * len_
    dfs(0, len_, div, mod, res)
    print(*res)
