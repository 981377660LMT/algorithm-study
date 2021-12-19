# 圆环上有10个点，编号为0~9。从0点出发，
# 每次可以逆时针和顺时针走一步，
# 问走n步回到0点共有多少种走法。

# 输入: 2
# 输出: 2
# 解释：有2种方案。分别是0->1->0和0->9->0
from functools import lru_cache


@lru_cache(None)
def dfs(cur: int, remain: int) -> int:
    if remain == 0:
        return int(cur == 0)
    return dfs((cur + 1) % 10, remain - 1) + dfs((cur - 1) % 10, remain - 1)


print(dfs(0, 2))

