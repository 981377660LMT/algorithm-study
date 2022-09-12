# 2 <= n <= 64
# 他想雨露均沾对待这两件事情，每件事情都刚好占用 n/2 天时间
# 他不想连续 d 天在做同一件事情(如果在某一件事情花费的时间已经到 n/2 天了，
# 剩下时间只能做另外一件事情，这种情况除外)
# 第一天的选择和最后一天的选择不一样
# 现在他想知道有多少种方案安排自己的时间


from functools import lru_cache


def solve():
    n, d = map(int, input().split())

    # 分别表示当前天数,上一个选择,上一个选择的连续天数,两个选择分别的剩余天数,以及第一天的选择
    @lru_cache(None)
    def dfs(index: int, pre: int, preLen: int, remain1: int, remain2: int, first: int) -> int:
        if index == n:
            return int(pre != first)

        # 只能做同一件事情
        if not remain1:
            return int(first == 0)
        if not remain2:
            return int(first == 1)

        res = 0
        # 能否继续上一个选择
        if preLen < d:
            if pre == 0:
                res += dfs(index + 1, 0, preLen + 1, remain1 - 1, remain2, first)
            else:
                res += dfs(index + 1, 1, preLen + 1, remain1, remain2 - 1, first)

        # 选择另外一个
        if pre == 0:
            res += dfs(index + 1, 1, 1, remain1, remain2 - 1, first)
        else:
            res += dfs(index + 1, 0, 1, remain1 - 1, remain2, first)
        return res

    res = dfs(1, 0, 1, n // 2 - 1, n // 2, 0) + dfs(1, 1, 1, n // 2, n // 2 - 1, 1)
    print(res)


if __name__ == "__main__":
    for _ in range(int(input())):
        solve()
