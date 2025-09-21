# E - Cut in Half (切成两半)
# https://atcoder.jp/contests/abc424/tasks/abc424_e
#
# 问题描述
#
# 一个袋子里有 N 根木棍，长度分别为 A_1, ..., A_N。
#
# 你将重复以下操作 K 次：
#
# 从袋子中取出最长的一根木棍。
# 将其对半切开。
# 将切开后的两根木棍放回袋中。
# 在 K 次操作结束后，袋中将有 N+K 根木棍。请找出其中第 X 长的木棍的长度。
#
# 你需要处理 T 个测试用例。
#
# 约束条件
#
# 1 ≤ T ≤ 10^5
# 对于每个测试用例：
# 1 ≤ N ≤ 10^5
# 1 ≤ A_i ≤ 10^9
# 1 ≤ K ≤ 10^9
# 1 ≤ X ≤ N+K
# 所有测试用例的 N 的总和不超过 10^5。
# 所有输入均为整数。


from math import ceil, floor, log2
from typing import Callable


def bisectRightFloat(
    left: float, right: float, check: Callable[[float], bool], absErrorInv=int(1e9)
) -> float:
    diff = ceil((right - left) * absErrorInv)
    round = diff.bit_length()
    for _ in range(round):
        mid = (left + right) / 2
        if check(mid):
            left = mid
        else:
            right = mid
    return (left + right) / 2


if __name__ == "__main__":

    def solve():
        _, K, X = map(int, input().split())
        A = list(map(int, input().split()))

        def check(mid: float) -> bool:
            c1, c2 = 0, 0
            for v in A:
                if v < mid:
                    continue
                c1 += 1
                c2 += (1 << floor(log2(v / mid))) - 1
            if K <= c2:
                return c1 + K >= X
            return c1 + c2 - (K - c2) >= X

        res = bisectRightFloat(0, max(A), check)
        print(res)

    T = int(input())
    for _ in range(T):
        solve()
