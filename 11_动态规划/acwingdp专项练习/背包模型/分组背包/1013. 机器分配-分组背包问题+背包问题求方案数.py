# 总公司拥有M台 相同 的高效设备，准备分给下属的N个分公司。
# 各分公司若获得这些设备，可以为国家提供一定的盈利。盈利与分配的设备数量有关。
# 问：如何分配这M台设备才能使国家得到的盈利最大？
# 求出最大盈利值。
# 分配原则：每个公司有权获得任意数目的设备，但总台数不超过设备数M。
# 1≤N≤10 ,
# 1≤M≤15

# 第一行有两个数，第一个数是分公司数N，第二个数是设备台数M；
# 接下来是一个N*M的矩阵，矩阵中的第 i 行第 j 列的整数表示第 i 个公司分配 j 台机器时的盈利。

# 输出格式
# 第一行输出最大盈利值；
# 接下N行，每行有2个数，即分公司编号和该分公司获得设备台数。

# 分组背包问题+背包问题求方案数

"""
每个公司看个一个物品组  取几个表示体积  对应盈利表示价值
"""
goodsCount, cap = map(int, input().split())
profits = []
for _ in range(goodsCount):
    # 矩阵中的第 i 行第 j 列的整数表示第 i 个公司分配 j 台机器时的盈利。
    # 注意要加上选0个的情况
    profits.append([0] + list(map(int, input().split())))

dp = [0] * (cap + 1)  # 最大盈利，和对应的选择个数
assign = [[0] * goodsCount for _ in range(cap + 1)]  # 每个容量下每个物品的选择个数

# 物品
for i in range(goodsCount):
    # 容量
    for j in range(cap, -1, -1):
        # 决策
        for k in range(j, -1, -1):
            if dp[j - k] + profits[i][k] > dp[j]:
                dp[j] = dp[j - k] + profits[i][k]
                # 求每个容量下的方案
                assign[j] = assign[j - k][:]
                assign[j][i] = k


print(dp[-1])
for i in range(goodsCount):
    print(i + 1, assign[cap][i])
