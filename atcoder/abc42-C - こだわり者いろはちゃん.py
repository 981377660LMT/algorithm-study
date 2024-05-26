# abc42-C - こだわり者いろはちゃん
# 求>=n的，不包含ban中数字的最小数
# !从个位开始考虑

N, K = map(int, input().split())
BAN = set(int(x) for x in input().split())


safe = set(range(10)) - BAN


def check(v: int) -> bool:
    return set(int(x) for x in str(v)) <= safe


def cal(n: int) -> int:
    while n % 10:
        if check(n):
            return n
        n += 1
    x = cal(n // 10)
    y = min(safe)
    return 10 * x + y


print(cal(N))
