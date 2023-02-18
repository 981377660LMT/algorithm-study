# quotient_range
# 数论分块


from typing import List, Tuple


def floorRange(n: int) -> List[Tuple[int, int, int]]:
    """
    将 [1,n] 内的数分成O(2*sqrt(n))段, 每段内的 n//i 相同

    Args:
        n (int): n>=1

    Returns:
        List[Tuple[int,int,int]]:
        每个元素为(left,right,div)
        表示 left <= i <= right 内的 n//i == div
    """
    res = []
    m = 1
    while m * m <= n:
        res.append((m, m, n // m))
        m += 1
    for i in range(m, 0, -1):
        left = n // (i + 1) + 1
        right = n // i
        if left <= right and res and res[-1][1] < left:
            res.append((left, right, n // left))
    return res


if __name__ == "__main__":
    n = int(input())
    print(floorRange(n))
    # [(1, 2, 9), (2, 3, 4), (3, 4, 3), (5, 10, 1)]
