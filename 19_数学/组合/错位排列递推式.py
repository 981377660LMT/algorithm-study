from functools import lru_cache


@lru_cache(None)
def cal(n: int) -> int:
    """错位排列递推式"""
    if n == 1:
        return 0
    if n == 2:
        return 1
    return (n - 1) * (cal(n - 1) + cal(n - 2))


print(cal(3))

