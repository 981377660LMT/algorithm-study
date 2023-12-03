from collections import deque
from typing import List


def toRootedTree(tree: List[List[int]], root=0) -> List[List[int]]:
    """无根树转有根树."""
    n = len(tree)
    rootedTree = [[] for _ in range(n)]
    visited = [False] * n
    visited[root] = True
    queue = deque([root])
    while queue:
        cur = queue.popleft()
        for next_ in tree[cur]:
            if not visited[next_]:
                visited[next_] = True
                queue.append(next_)
                rootedTree[cur].append(next_)
    return rootedTree


if __name__ == "__main__":
    tree = [[1, 2], [0, 3, 4], [0, 5, 6], [1], [1], [2], [2]]
    print(toRootedTree(tree))
