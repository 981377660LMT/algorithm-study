# 指定元素,生成大小为r的置换环

from itertools import permutations
from typing import Generator, List, Tuple, TypeVar

T = TypeVar("T")


def genCycle(elements: List[T], r: int) -> Generator[Tuple[T, ...], None, None]:
    for perm in permutations(elements, r):
        perm = perm + (perm[0],)
        yield perm


if __name__ == "__main__":
    # 例如元素为1,2,3,4 大小为2的置换环:
    elements = [1, 2, 3, 4]
    r = 2
    for cycle in genCycle(elements, r):
        for a, b in zip(cycle, cycle[1:]):
            print(a, "->", b, end=" ")
        print()
