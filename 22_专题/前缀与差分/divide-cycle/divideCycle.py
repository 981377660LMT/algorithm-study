# 环区间分解


from random import randint
from typing import Generator, Tuple


def divideCycle(n: int, start: int, end: int) -> Generator[Tuple[int, int, int], None, None]:
    """将环上的区间分解为形如`[start, end)`的区间，每个区间遍历`times`次."""
    if start >= end or n <= 0:
        return
    loop = (end - start) // n
    if loop > 0:
        yield 0, n, loop
    if (end - start) % n == 0:
        return
    start %= n
    end %= n
    if start < end:
        yield start, end, 1
    else:
        yield start, n, 1
        if end > 0:
            yield 0, end, 1


if __name__ == "__main__":
    for start, end, times in divideCycle(10, 8, 20):
        print(start, end, times)

    arr = [randint(0, 100) for _ in range(100)]

    def bf(n, start, end) -> int:
        res = 0
        for i in range(start, end):
            res += arr[i % n]
        return res

    def calc(n, start, end):
        res = 0
        for s, e, t in divideCycle(n, start, end):
            res += t * sum(arr[s:e])
        return res

    for _ in range(1000):
        n = randint(1, 100)
        start, end = randint(0, n * 8), randint(0, n * 8)
        if start > end:
            start, end = end, start
        if bf(100, start, end) != calc(100, start, end):
            print("error")
            break
    print("ok")
