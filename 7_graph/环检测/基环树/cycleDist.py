# 给定一个环，环的节点编号为0~n-1.


from typing import List


def distOnCycle(n: int, i: int, j: int) -> int:
    """环上两点距离."""
    cand1, cand2 = abs(i - j), n - abs(i - j)
    return cand1 if cand1 < cand2 else cand2


def distPairOnCycle(n: int) -> List[int]:
    """环上两点距离为k的点对数."""
    res = [0] * (n // 2 + 1)
    for k in range(1, n // 2 + 1):
        res[k] = n
    if n % 2 == 0:
        res[n // 2] = n // 2
    return res


if __name__ == "__main__":
    print(distPairOnCycle(3))  # [0,3]
    print(distPairOnCycle(4))  # [0,4,2]
    print(distPairOnCycle(5))  # [0,5,5]
