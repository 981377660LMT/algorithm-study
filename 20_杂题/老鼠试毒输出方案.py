# 老鼠试毒输出方案
# https://atcoder.jp/contests/abc337/editorial/9140

# 交互题
# !1.结论需要的老鼠数need满足 2**need>=n
# !2.将n转为二进制，第i只老鼠需要喝二进制下第i位为1的毒药
# !3.如果第i只老鼠死亡，则第i位为1，否则为0


from typing import Callable, List, Tuple


def solve(n: int) -> Tuple[List[List[int]], Callable[[List[bool]], int]]:
    """返回老鼠试毒方案.n瓶毒药中有且仅有一瓶有毒,老鼠试毒后可以确定哪瓶有毒.

    Args:
        n (int): 毒药数

    Returns:
        Tuple[List[List[int]], Callable[[List[bool]], int]]:
        老鼠试毒方案, 根据老鼠试毒结果解析出有毒毒药编号.
    """
    need = 0
    while (1 << need) < n:
        need += 1

    res = []
    for i in range(need):
        cur = [j for j in range(n) if j & (1 << i)]
        res.append(cur)

    def parse(dead: List[bool]) -> int:
        return sum(1 << i for i, v in enumerate(dead) if v)

    return res, parse


if __name__ == "__main__":
    n = int(input())

    res, parse = solve(n)

    print(len(res), flush=True)
    for row in res:
        print(len(row), *[v + 1 for v in row], flush=True)

    s = input()
    dead = [c == "1" for c in s]
    print(parse(dead) + 1, flush=True)
