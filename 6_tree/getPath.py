# 求树中两点路径


from typing import List


def getPath(tree: List[List[int]], from_: int, to: int) -> List[int]:
    def dfs(cur: int, pre: int) -> bool:
        path.append(cur)
        if cur == to:
            return True
        for next in tree[cur]:
            if next != pre and dfs(next, cur):
                return True
        path.pop()
        return False

    path = []
    dfs(from_, -1)
    return path


if __name__ == "__main__":
    n = 5
    tree = [[] for _ in range(n)]
    edges = [[0, 1], [0, 2], [1, 3], [1, 4]]
    for u, v in edges:
        tree[u].append(v)
        tree[v].append(u)
    print(getPath(tree, 0, 4))
