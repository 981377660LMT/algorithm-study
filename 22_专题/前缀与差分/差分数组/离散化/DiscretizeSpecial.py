from typing import List, Tuple


def discretizeSpecial(n: int, isSpecial: List[bool]) -> Tuple[List[int], List[int]]:
    """Discretize arr to id array.

    >>> discretizeSpecial(5, [True, False, True, False, True])
    ([0, -1, 1, -1, 2], [0, 2, 4])
    """
    vToId = [-1] * n
    idToV = []
    for i in range(n):
        if isSpecial[i]:
            vToId[i] = len(idToV)
            idToV.append(i)
    return vToId, idToV


if __name__ == "__main__":
    assert discretizeSpecial(5, [True, False, True, False, True]) == (
        [0, -1, 1, -1, 2],
        [0, 2, 4],
    )
