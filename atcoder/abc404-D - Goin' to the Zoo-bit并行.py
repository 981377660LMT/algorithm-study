# !三进制枚举
# https://atcoder.jp/contests/abc404/editorial/12874

INF = int(1e18)
if __name__ == "__main__":
    n, m = map(int, input().split())
    cost = tuple(map(int, input().split()))
    zoo = [[] for _ in range(n)]
    for animal in range(m):
        _, *ids = map(lambda x: int(x) - 1, input().split())
        for j in ids:
            zoo[j].append(animal)

    ones = [sum(1 << 2 * j for j in z) for z in zoo]
    twoAll = sum(2 << 2 * j for j in range(m))

    def addOne(watched: int, one: int) -> int:
        return watched + (one & ~watched >> 1)

    def dfs(i: int, watched: int, res: int) -> int:
        if i == n:
            return res if watched == twoAll else INF
        res1 = dfs(i + 1, watched, res)
        watched = addOne(watched, ones[i])
        res2 = dfs(i + 1, watched, res + cost[i])
        watched = addOne(watched, ones[i])
        res3 = dfs(i + 1, watched, res + cost[i] * 2)
        return min(res1, res2, res3)

    print(dfs(0, 0, 0))
