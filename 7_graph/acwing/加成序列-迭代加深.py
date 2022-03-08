def dfs(cur: int, limit: int) -> bool:
    if cur > limit:
        return False
    if path[cur - 1] == t:
        return True
    flag = [False] * N
    for i in range(cur - 1, -1, -1):
        for j in range(i, -1, -1):
            y = path[i] + path[j]
            if y > t or y <= path[cur - 1] or flag[y]:
                continue
            flag[y] = True
            path[cur] = y
            if dfs(cur + 1, limit):
                return True
    return False


N = 110
while True:
    t = int(input())
    if not t:
        break
    path = [0] * N
    path[0] = 1
    depth = 1
    while not dfs(1, depth):
        depth += 1
    print(*path[:depth])

