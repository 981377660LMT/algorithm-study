from functools import lru_cache

MOD = int(1e9 + 7)


res = [0, 0, 1]
for i in range(3, int(1e6) + 10):
    res.append(((i - 1) * (res[-1] + res[-2])) % MOD)

###########################################################
@lru_cache(None)
def cal(n: int) -> int:
    """错位排列递推式"""
    if n == 1:
        return 0
    if n == 2:
        return 1
    return (n - 1) * (cal(n - 1) + cal(n - 2))
