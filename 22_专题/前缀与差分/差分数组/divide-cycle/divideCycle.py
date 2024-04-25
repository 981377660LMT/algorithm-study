# 环区间分解


from typing import Generator, Tuple


def divideCycle(n: int, start: int, end: int) -> Generator[Tuple[int, int, int], None, None]:
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
