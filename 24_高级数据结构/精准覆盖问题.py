# 精准覆盖问题


from typing import List


def exactOver(sets: List[List[int]]) -> List[int]:
    def remove(x: int) -> None:
        ...

    def resume(x: int) -> None:
        ...

    def rec() -> bool:
        ...

    m = 0
    M = 10
    for s in sets:
        m = max(s, default=0) + 1
    M += m * (1 + len(sets))
    solution = []
    rec()
    return solution
