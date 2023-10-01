# CyclePartition-置换环分组+容斥原理


from typing import List

MOD = 998244353


def cyclePartition(nexts: List[int]) -> List[List[int]]:
    """给定一个0-n-1的排列,返回环分组
    nexts[i]表示i的下一个元素.
    """
    ...


def dontBeTogether(perm: List[int], m: int) -> int:
    ...


if __name__ == "__main__":
    # https://yukicoder.me/problems/5125
    # 将1 - n的整数分为m组, 且i和perm[i]不能一组
    # 求方案数

    # 容斥原理:
    # 求i和perm[i]在一组的方案数

    n, m = map(int, input().split())
    perm = [int(x) - 1 for x in input().split()]
    print(dontBeTogether(perm, m))

    # TODO
    # https://yukicoder.me/submissions/617293
