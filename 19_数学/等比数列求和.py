# 求前n项和
from typing import List


# 等比数列求和
def getSum(n: int, a0: int, q: int) -> List[int]:
    """等比数列前n项和"""
    res = [a0]
    curSum, curItem = a0, a0
    for _ in range(n - 1):
        curItem *= q
        curSum += curItem
        res.append(curSum)
    return res


print(getSum(5, 2, 2))
