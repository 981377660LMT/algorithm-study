from typing import List


def getFactorsAll(max: int) -> List[List[int]]:
    """预处理 1~max 的所有数的因数."""
    res = [[] for _ in range(max + 1)]
    for f in range(1, max + 1):
        for m in range(f, max + 1, f):
            res[m].append(f)
    return res


assert getFactorsAll(10) == [  # noqa
    [],
    [1],
    [1, 2],
    [1, 3],
    [1, 2, 4],
    [1, 5],
    [1, 2, 3, 6],
    [1, 7],
    [1, 2, 4, 8],
    [1, 3, 9],
    [1, 2, 5, 10],
]
