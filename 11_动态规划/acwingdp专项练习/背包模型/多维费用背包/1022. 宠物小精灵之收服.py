'''
本质是01背包，约束条件有两个维度，价值也是有两个维度（收服的怪兽数目，剩余血量），需要把
价值映射成1维度的来比较大小进行决策


dp(i, j, k) 表示前i种怪兽在精灵球剩余j个，血量剩余k条件下的最大价值
实现的时候把第一维的空间压缩了
'''
# 输入数据的第一行包含三个整数：N，M，K，分别代表小智的精灵球数量、皮卡丘初始的体力值、野生小精灵的数量。
# 之后的K行，每一行代表一个野生小精灵，包括两个整数：收服该小精灵需要的精灵球的数量，以及收服过程中对皮卡丘造成的伤害。
# 输出为一行，包含两个整数：C，R，分别表示最多收服C个小精灵，以及收服C个小精灵时皮卡丘的剩余体力值最多为R。
# 0<N≤1000,
# 0<M≤500,
# 0<K≤100
# 小智遇到野生小精灵时有两个选择：收服它，或者离开它。  (要不要选择


ball, hp, monster = map(int, input().split())
goods = []
for _ in range(monster):
    ballCost, damage = map(int, input().split())
    goods.append((ballCost, damage))

# 二维费用
dp = [[0] * (hp + 1) for _ in range(ball + 1)]
for ballCost, damage in goods:
    for i in range(ball, ballCost - 1, -1):
        for j in range(hp, damage, -1):
            # 血量不能为0
            if i >= ballCost - 1 and j >= damage + 1:
                dp[i][j] = max(dp[i][j], dp[i - ballCost][j - damage] + 1)

threshold = next((i for i in range(hp, -1, -1) if dp[ball][i] != dp[ball][hp]), 0)
print(dp[ball][hp], hp - threshold)

