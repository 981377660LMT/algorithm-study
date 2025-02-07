# E - Hierarchical Majority Vote(三叉树树形dp)
# https://atcoder.jp/contests/abc391/editorial/12106
# 给定一个长为3^N 的二进制字符串 A，通过 N 次操作将其压缩为长度为 1 的字符串 A1′。
# 每次操作将字符串分成长度为 3 的组，取每组的多数值作为新字符串的元素。
# 目标是通过最少的修改次数（将 0 改为 1 或将 1 改为 0），使得最终的 A1′ 发生变化。
# !即：我们最终希望通过最少的位变化，将经过N次操作后得到的字符串的结果从当前值变为相反的值。
# 类似力扣 1896. 反转表达式的最少操作次数

from typing import Tuple


def hierarchicalMajorityVote(N: int, S: str) -> int:
    def dfs(level: int, offset: int) -> Tuple[int, int]:
        """变为0/1的最小反转次数."""
        if level == -1:
            return (0, 1) if S[offset] == "0" else (1, 0)
        zero, one = zip(*[dfs(level - 1, offset + i * 3**level) for i in range(3)])
        res0 = sum(zero) - max(zero)
        res1 = sum(one) - max(one)
        return res0, res1

    res = dfs(N - 1, 0)
    return max(res)


if __name__ == "__main__":
    N = int(input())
    S = input()
    print(hierarchicalMajorityVote(N, S))
