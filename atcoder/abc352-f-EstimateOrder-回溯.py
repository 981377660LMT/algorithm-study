# F - EstimateOrder
# https://atcoder.jp/contests/abc352/tasks/abc352_f
# n个人，排名唯一，给定关于这些人的m条限制，问每个人的名次是否能唯一确定。
# 每条排名消息形如： a - b = c，表示a比b名次高c名。
# n<=16.
# 如果人i的排名唯一确定，输出其排名，否则输出-1。
# 保证至少存在一种排名满足所有限制。
#
# !回溯法(状压dp不好写，考虑回溯爆搜)


import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    n, m = map(int, input().split())
    limits = []
    for _ in range(m):
        a, b, c = map(int, input().split())
        limits.append((a - 1, b - 1, c))

    rankToPerson = [-1] * n
    candidates = [[False] * n for _ in range(n)]
    assignCount = [0] * n  # 每个人被分配名次的次数

    def bt(index: int) -> None:
        if index == m:
            unusedRank = []
            person = [False] * n
            for r, p in enumerate(rankToPerson):
                if p == -1:
                    unusedRank.append(r)
                else:
                    person[p] = True
                    candidates[p][r] = True
            for p, b in enumerate(person):
                if not b:
                    for r in unusedRank:
                        candidates[p][r] = True
            return

        a, b, c = limits[index]
        for rb in range(n - c):
            ra = rb + c
            if (rankToPerson[ra] != -1 and rankToPerson[ra] != a) or (
                rankToPerson[rb] != -1 and rankToPerson[rb] != b
            ):
                continue
            preA, preB = rankToPerson[ra], rankToPerson[rb]
            assignCount[a] += rankToPerson[ra] == -1
            assignCount[b] += rankToPerson[rb] == -1
            rankToPerson[ra], rankToPerson[rb] = a, b
            # !合法的rank分配
            if assignCount[a] == assignCount[b] == 1:
                bt(index + 1)
            rankToPerson[ra], rankToPerson[rb] = preA, preB
            assignCount[a] -= rankToPerson[ra] == -1
            assignCount[b] -= rankToPerson[rb] == -1

    limits.sort(key=lambda x: -x[2])  # ! 剪枝, 优先分配名次差距最大的(c越大，可能的rank越少)
    bt(0)
    res = [-1] * n
    for i, row in enumerate(candidates):
        if row.count(True) == 1:
            res[i] = row.index(True) + 1
    print(" ".join(map(str, res)))
