# 有多少种？第 N 项 Catalan 数，即 C(2n,n)/(n +1)

# 并求出顺序

from typing import List


def dfs(index: int, inTrains: List[int], outTrains: List[int]) -> None:
    """当前index，已经入栈的火车，已经出栈的火车"""
    if trains[-1] in inTrains:
        res.append(' '.join(outTrains + inTrains[::-1]))
        return

    if not inTrains:
        dfs(index + 1, inTrains + [trains[index]], outTrains)
    else:
        dfs(index + 1, inTrains + [trains[index]], outTrains)
        dfs(index, inTrains[:-1], outTrains + [inTrains[-1]])


while True:
    try:
        n = int(input())
        # 火车入站的序列
        trains = input().split()
        res = []
        dfs(0, [], [])
        res.sort()
        print('\n'.join(res))
    except EOFError:
        break

