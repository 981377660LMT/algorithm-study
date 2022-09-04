from typing import List


def findLCA(tree: List[int], val1: int, val2: int) -> int:
    mp = {num: i for i, num in enumerate(tree)}
    if val1 not in mp or val2 not in mp:
        return -1
    i1, i2 = mp[val1], mp[val2]
    path = set([0])

    while i1:
        path.add(i1)
        i1 = (i1 - 1) // 2

    while i2:
        if i2 in path:
            return tree[i2]
        i2 = (i2 - 1) // 2
    return tree[0]


print(findLCA([5, 2, 4, 1, 6, 9, 0, 3], 3, 6))
