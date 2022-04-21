N = 12
dp = [[0] * (1 << N) for _ in range(N)]
validStates = [False] * (1 << N)


while True:
    n, m = map(int, input().split())
    if 0 in (n, m):
        break

    # 预处理一列的所有可能状态
    for state in range(1 << n):
        # cnt表示本列当前连续的空格子的数量
        count = 0
        validStates[state] = True
        # 一位一位的判断每个状态
        for i in range(n):
            if (state >> i) & 1:
                # cnt&1; 等价于cnt%2!=0，即当前有奇数个空格子，就不能竖着放满方块
                if count & 1:
                    validStates[state] = False
                count = 0
            else:
                count += 1

        # 遍历完所有后剩下奇数个格子，则该状态不合法
        if count & 1:
            validStates[state] = False

    # 因为是多组数据，每一组数据用之前都要把上一次的记录清空
    dp = [[0] * (1 << n + 1) for _ in range(m + 1)]

    dp[0][0] = 1
    # 枚举每一列
    for col in range(1, m + 1):
        # 枚举第i列的所有状态
        for state in range(1 << n):
            # 枚举第i-1列的所有状态
            for preState in range(1 << n):
                # 判断两列匹配是否合法
                if state & preState == 0 and validStates[state | preState]:
                    dp[col][state] += dp[col - 1][preState]
    print(dp[m][0])

# print(st)


# 太烦了 不看了 思路就是这样的
