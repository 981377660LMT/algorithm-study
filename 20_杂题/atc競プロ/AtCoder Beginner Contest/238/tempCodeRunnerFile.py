    def dfs(index: int, remain: int, preMin: int) -> int:
        if remain < 0:
            return 0
        if index == n:
            return 1 if remain == 0 else 0
        hash = index * h2 + remain * h1 + preMin
        if memo[hash] != -1:
            return memo[hash]

        res = dfs(index + 1, remain, min(preMin, people[index][1]))  # jump
        if remain > 0 and people[index][1] < preMin:
            res += dfs(index + 1, remain - 1, preMin)
        return res % MOD

    n, k = map(int, input().split())
    rank1 = list(map(int, input().split()))
    rank2 = list(map(int, input().split()))
    people = sorted(zip(rank1, rank2), key=lambda x: x[0])
    h3, h2, h1 = (n + 5) * (n + 5) * (n + 5), (n + 5) * (n + 5), (n + 5)
    memo = [-1] * h3
    print(dfs(0, k, n + 1))