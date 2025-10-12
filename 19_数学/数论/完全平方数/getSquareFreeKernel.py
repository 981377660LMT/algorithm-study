from typing import Callable, Mapping


def getSquareFreeKernel(n: int, getPrimeFactors: Callable[[int], Mapping[int, int]]) -> int:
    res = 1
    for p, c in getPrimeFactors(n).items():
        if c & 1:
            res *= p
    return res
