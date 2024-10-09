# 快速因数分解


from typing import List, Tuple


def getPrimeFactors(n: int) -> List[Tuple[int, int]]:
    """质因数分解.

    >>> getPrimeFactors(100)
    [(2, 2), (5, 2)]
    """
    res = []
    upper = n
    i = 2
    while i * i <= upper:
        if upper % i == 0:
            c = 0
            while upper % i == 0:
                c += 1
                upper //= i
            res.append((i, c))
        i += 1
    if upper != 1:
        res.append((upper, 1))
    return res


if __name__ == "__main__":
    print(getPrimeFactors(int(1e16 + 91)))
