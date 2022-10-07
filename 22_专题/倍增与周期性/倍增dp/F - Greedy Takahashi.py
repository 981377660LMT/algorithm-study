"""城市公交漫游/开车旅行

# https://www.cnblogs.com/Acfboy/p/ABC212F.html
# n个城市,m辆车,对于第i辆车在si+0.5时刻从u出发,于ti+0.5抵达vi。
# 所有的公交车出发时间不同.
# T某在一个点时会选择还没出发且离当前出发时刻最近的车坐上,直到没有车。
# !给出q个询问,时刻Xi时T某从Yi城市出发,问Zi时刻T某在哪个城市。
# !n,m,q<=1e5
# !si,ti<=1e9

倍增加速模拟(注意到每次转移都具有唯一性)
dp[i][j] 表示乘坐编号为i的车后又换乘2^j辆车后乘坐的车的编号(从一辆车转移到另一辆车)
dp[i][0] 可以用二分查找求出
"""


from math import floor, log2
import sys
from sortedcontainers import SortedList

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, q = map(int, input().split())
    bus = []  # !(city1, city2, startTime, endTime)
    city = [SortedList() for _ in range(n)]  # !每个城市出发的公交车 (startTime, i)
    for i in range(m):
        city1, city2, startTime, endTime = map(int, input().split())
        city1, city2 = city1 - 1, city2 - 1
        bus.append((city1, city2, startTime, endTime))
        city[city1].add((startTime, i))

    maxJ = floor(log2(m)) + 1
    dp = [[-1] * (maxJ + 1) for _ in range(m + 1)]  # !dp[i][j] 表示乘坐编号为i的车后又换乘2^j辆车后乘坐的车的编号
    for i in range(m):
        city1, city2, startTime, endTime = bus[i]
        pos = city[city2].bisect_left((endTime, -1))
        if pos < len(city[city2]):
            dp[i][0] = city[city2][pos][1]

    for j in range(1, maxJ + 1):
        for i in range(m):
            if dp[i][j - 1] != -1:
                dp[i][j] = dp[dp[i][j - 1]][j - 1]

    for _ in range(q):
        startTime, startCity, queryTime = map(int, input().split())
        startCity -= 1
        pos = city[startCity].bisect_left((startTime, -1))
        if pos == len(city[startCity]) or bus[city[startCity][pos][1]][2] >= queryTime:
            print(startCity + 1)  # !没有车出发或者最早的车出发时间晚于queryTime
        else:
            curBus = city[startCity][pos][1]
            for j in range(maxJ, -1, -1):
                if dp[curBus][j] != -1 and bus[dp[curBus][j]][2] < queryTime:
                    curBus = dp[curBus][j]

            # 在坐巴士
            if bus[curBus][3] >= queryTime:
                city1, city2 = bus[curBus][0], bus[curBus][1]
                print(city1 + 1, city2 + 1)
            # 到达了城市
            else:
                city2 = bus[curBus][1]
                print(city2 + 1)
