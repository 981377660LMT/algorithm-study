# 小美最近在玩一款回合制游戏。该游戏中的一个人机关卡共有n个回台。
# 小美所使用的角色并没有主动攻击技能,有且只有一个防御技能——圣光!
# 圣光总共可以使用k次,每回最多使用一次,如果在某一个回合使用了圣光.
# 那么小美的角色会免疫该回合受到的所有伤害，
# !并且在之后没有释放圣光的回合中每回合恢复1点HP,效果可以叠加，
# 例如在前两个回合都释放了一次圣光,那么第三个回合就会恢复2点HP.
# 小美只需要在n个回合后HP大于等于0，即可通关，
# 请问小美角色初始最少需要多少HP.
# (初始HP应该大于等于1,游戏过程中允许部分时刻HP小于0)
# n,k,damage[i]<=1e4


# 1e4
# !dp[index][use] 表示在第index个回合使用use次圣光时的最小损伤
# use or not use
# dp[index]=min(dp[])
# from collections import Counter


# INF = int(1e10)
# n, k = map(int, input().split())  # 回台数，圣光次数
# attack = list(map(int, input().split()))  # 每个回台的伤害
# dp = [0] * (k + 1)
# for i in range(n):
#     ndp = [INF] * (k + 1)
#     for j in range(k + 1):
#         # 不用圣光
#         ndp[j] = min(ndp[j], dp[j] + attack[i] - j)
#         # 用圣光
#         if j + 1 <= k:
#             ndp[j + 1] = min(ndp[j + 1], dp[j])
#     dp = ndp

# minDamage = min(dp)
# print(max(1, minDamage))
