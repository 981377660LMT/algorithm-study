from itertools import chain, combinations
from typing import Collection, List, Tuple


def powerset(collection: Collection, includeAll=True):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) + 1 if includeAll else len(collection)
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


def genPowerSetFromAllPowerSet(collection: Collection):
    """枚举所有子集的子集，时间复杂度O(3^n)"""
    allPowetSet = powerset(collection)
    return chain.from_iterable(powerset(subset) for subset in allPowetSet)


def genPowerSetFromAllPowerSet2(nums: List[int]) -> List[List[Tuple[int, ...]]]:
    """举所有子集的子集，返回pair互补对，时间复杂度O(3^n)"""
    n = len(nums)
    res = []
    for state in range(1 << n):
        cur = []
        g1, g2 = state, 0
        while g1:  # 如果子集规定非空，这里变成while group1 即可
            cur.append((g1, g2))  # 其实append group1就可以了
            # 关键，不断减一+与运算跳数
            g1 = state & (g1 - 1)
            g2 = state ^ g1
        res.append(cur)
    return res


def genPowerSetFromAllPowerSet3(n: int):
    """枚举每个长度不超过n的二进制数的子集，时间复杂度O(3^n)"""
    res = []
    for state in range(1 << n):
        cur = []
        g1 = state
        while g1:
            cur.append(g1)
            g1 = state & (g1 - 1)
        res.append(cur)
    return res


if __name__ == "__main__":
    print(len(list(genPowerSetFromAllPowerSet([1, 2, 3, 4]))))
    print(len(sum(genPowerSetFromAllPowerSet2([1, 2, 3, 4]), [])))
    print(len(sum(genPowerSetFromAllPowerSet3(4), [])))
