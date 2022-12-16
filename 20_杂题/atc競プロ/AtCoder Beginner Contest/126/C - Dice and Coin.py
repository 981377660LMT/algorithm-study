# C - Dice and Coin
# 骰子与硬币

# Snuke 有一个 n面的色子，投掷这个色子的时候会以相等的概率得到一个在 1到 n 之间的整数。
# 他还有一个硬币，投掷时正面朝上和反面朝上的概率相等。

# 现在他要用色子和硬币玩一个游戏
# 扔色子，将得到的整数作为初始分数。
# 1.只要这个分数在 1 到 k−1 之间（包含 1 和 k−1），就扔硬币。
# 当正面朝上时，将这一分数翻倍；否则，将分数归零。
# 2.分数归零或大于等于 k 时，游戏结束。若分数大于等于 k，Snuke 获胜，否则 Snuke 失败。
# 给出 n 和 k，你需要求出 Snuke 获胜的概率。
# n,k<=1e5
# !模拟即可
def diceAndCoin(n: int, k: int) -> float:
    res = 0
    for dice in range(1, n + 1):
        score = dice
        cur = 1 / n
        while score < k:  # 需要赢
            score *= 2
            cur /= 2
        res += cur
    return res


n, k = map(int, input().split())
print(diceAndCoin(n, k))
