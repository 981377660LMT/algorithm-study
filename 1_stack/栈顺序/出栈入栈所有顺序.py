# 有多少种？第 N 项 Catalan 数，即 C(2n,n)/(n +1)

# 并求出顺序

# // 面对任何一个状态, 我们只有两种选择:
# // 1. 把下一个数进栈(如果还有下一个数)
# // 2. 把当前栈顶的数出栈(如果栈非空)

from typing import List


def dfs(index: int, inTrains: List[int], outTrains: List[int]) -> None:
    """当前index，已经入栈的火车，已经出栈的火车"""

    if len(outTrains) == n:
        res.append(outTrains[:])
        return

    # 1. 下一个数进站
    if index <= n:
        inTrains.append(index)
        dfs(index + 1, inTrains, outTrains)
        inTrains.pop()

    # 2. 栈顶出栈
    if inTrains:
        outTrains.append(inTrains.pop())
        dfs(index, inTrains, outTrains)
        inTrains.append(outTrains.pop())


while True:
    try:
        n = int(input())
        res = []
        dfs(1, [], [])
        res.sort()
        print(res)
    except EOFError:
        break

