from typing import List


def restorePath(target: int, pre: List[int]) -> List[int]:
    """还原路径/dp复原."""
    path = [target]
    while pre[path[-1]] != -1:
        path.append(pre[path[-1]])
    path.reverse()
    return path
