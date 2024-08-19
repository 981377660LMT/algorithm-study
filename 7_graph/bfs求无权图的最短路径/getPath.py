from typing import List


def getPath(graph: List[List[int]], start: int, target: int) -> List[int]:
    n = len(graph)
    visited = [False] * n
    pre = [-1] * n
    stack = [start]
    visited[start] = True
    while stack:
        cur = stack.pop()
        if cur == target:
            break
        for next_ in graph[cur]:
            if not visited[next_]:
                visited[next_] = True
                pre[next_] = cur
                stack.append(next_)
    if not visited[target]:
        return []
    res = []
    while target != -1:
        res.append(target)
        target = pre[target]
    return res[::-1]


if __name__ == "__main__":
    print(getPath([[1, 2], [3], [3], []], 0, 3))  # [0, 2, 3]
    print(getPath([[1, 2], [3], [3], []], 0, 1))  # [0, 1]
    print(getPath([[1, 2], [3], [3], []], 0, 0))  # [0]
    print(getPath([[1, 2], [3], [3], []], 0, 2))  # [0, 2]
