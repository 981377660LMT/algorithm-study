# 替换括号序列中的问号(括号序列,问号替换)
# 给一个序列，序列里面会有左括号、问号、右括号。对于一个‘？’而言，
# 可以将其替换为一个‘（’，也可以替换成一个‘）’，但是都有相应的代价。问：如何替换使得代价最小。
# 前提是替换之后的序列中，括号是匹配的。如果不能替换为一个括号匹配的序列则输出-1。
#
# 第一行输出最小代价，第二行输出替换后的序列。不行输出 −1。


# 思路:
# 对于每一个 ? 我们先把它当做 ) 丢进原序列。如果发现某个位置不合法就把从 ) 修改到 ( 代价差最小的进行更换。
# 把一个 ? 修改的时候只要丢到一个小根堆里面就好了


from heapq import heappop, heappush
from typing import List, Optional, Tuple


def leastCostBracketSequence(
    s: str, cost1: List[int], cost2: List[int]
) -> Optional[Tuple[int, str]]:
    diff = 0
    sb = list(s)
    res = 0
    pq = []

    for i, c in enumerate(s):
        if c == "(":
            diff += 1
        elif c == ")":
            diff -= 1
        else:
            diff -= 1
            sb[i] = ")"
            heappush(pq, (cost1[i] - cost2[i], i))
            res += cost2[i]
        if diff < 0 and not pq:  # 栈空
            return None
        if diff < 0:
            costDiff, index = heappop(pq)
            diff += 2
            res += costDiff
            sb[index] = "("
    if diff != 0:  # 非空栈
        return None
    return res, "".join(sb)


if __name__ == "__main__":
    s = input().strip()
    replaceCost1 = [-1] * len(s)
    replaceCost2 = [-1] * len(s)
    for i, c in enumerate(s):
        if c == "?":
            c1, c2 = map(int, input().split())
            replaceCost1[i] = c1
            replaceCost2[i] = c2

    res = leastCostBracketSequence(s, replaceCost1, replaceCost2)
    if res is None:
        print(-1)
    else:
        print(res[0])
        print(res[1])
