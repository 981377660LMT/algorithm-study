# 覆盖区间的方案数/区间覆盖的方案数
# 给定[0,X]内的n个区间
# 选择其中若干个,使得这些区间能够覆盖[0,X]内的所有点
# 求方案数
# !n<=1e3

# !按终点排序,dp[i][mex]表示前i个区间,未被覆盖的最左端点为mex的方案数
# 选区间的问题都可以用数据结构优化
#  把区间按照右端点排序 dp[i][mex] 然后用区间更新
# !注意dp优化要用貰うdp


from typing import List, Tuple

MOD = int(1e9 + 7)


def solve(intervals: List[Tuple[int, int]], x: int) -> int:
    """按照区间右端点排序,dp[i][mex]表示前i个区间,未被覆盖的最左端点为mex的方案数."""
    intervals.sort(key=lambda v: v[1])
    maxEnd = max(x, intervals[-1][1])
    dp = [0] * (maxEnd + 2)
    dp[0] = 1
    for s, e in intervals:
        ndp = dp[:]
        for curMex in range(s):
            ndp[curMex] += dp[curMex]
            ndp[curMex] %= MOD
        ndp[e + 1] += sum(dp[s : e + 2])
        ndp[e + 1] %= MOD
        dp = ndp
    return sum(dp[x + 1 :])


if __name__ == "__main__":

    def bruteForce(intervals: List[Tuple[int, int]], x: int) -> int:
        """暴力枚举,用于验证正确性."""
        n = len(intervals)
        res = 0
        for s in range(1, 1 << n):
            cur = set()
            for i in range(n):
                if s >> i & 1:
                    cur |= set(range(intervals[i][0], intervals[i][1] + 1))
            if all(v in cur for v in range(x + 1)):
                res += 1
        return res

    from random import randint

    for _ in range(100):
        # rando
        n = randint(1, 12)
        x = randint(1, 1000)
        intervals = []
        for _ in range(n):
            s = randint(0, x + 100)
            e = randint(s, x + 100)
            intervals.append((s, e))
        if solve(intervals, x) != bruteForce(intervals, x):
            print(intervals, x)
            print(solve(intervals, x), bruteForce(intervals, x))
            exit(0)
