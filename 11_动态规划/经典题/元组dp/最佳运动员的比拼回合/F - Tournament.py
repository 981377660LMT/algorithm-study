# 共有2"个人，他们排列在一排，两两进行比赛，2×i＋1和2×i
# 胜者进行下一轮继续，这样一共进行k轮，对于每个参赛者，
# !他们连胜若干场会有一个奖励，考虑如何安排胜负使每个参赛者奖励的和最大化


# 最佳运动员的比拼回合3

# !dfs(root,win) 表示root结点的赢家连续赢了win回合时赢家的奖励分数 (根节点为1)

n = int(input())
score = []
for _ in range(1 << n):
    row = list(map(int, input().split()))
    score.append([0] + row)


def dfs(root: int, win: int) -> int:
    if root >= (1 << n):  # 叶子结点表示赢家
        winner = root ^ (1 << n)
        return score[winner][win]
    hash_ = root * (n + 1) + win
    if ~memo[hash_]:
        return memo[hash_]
    res = dfs(root << 1, 0) + dfs(root << 1 | 1, win + 1)
    res = max(res, dfs(root << 1, win + 1) + dfs(root << 1 | 1, 0))
    memo[hash_] = res
    return res


memo = [-1] * (n + 5) * (1 << (n + 1))
print(dfs(1, 0))
