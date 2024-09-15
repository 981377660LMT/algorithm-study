from typing import List


def collectCycle(nexts: List[int], start: int) -> List[int]:
    """置换环找环.nexts数组中元素各不相同."""
    cycle = []
    cur = start
    while True:
        cycle.append(cur)
        cur = nexts[cur]
        if cur == start:
            break
    return cycle
