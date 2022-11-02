# f-捕鱼达人 扫描线
# 有n条鱼在数轴上移动。
# 第i条鱼在时刻0时在位置xi处，价值为wi，将会以每时刻ti的速度向数轴正方向前进。
# 你是一个渔夫，你有感应河流的能力，你已经知晓所有鱼的x, w,t属性。
# 你会选择一个时刻t，在位置x撒下一张长度为a的网，所有在时刻t时处于区间[x, x + a]的鱼都会被你捕获。
# !你想求出你撒一次网能捕获的鱼的价值和的最大值。
# !n <= 2e3 a,wi,xi,ti <= 1e4

# 枚举每一条鱼，以这条鱼为参照物，
# 假设他在捕鱼网内的最左端，
# 其他鱼于他的距离绝对值<=a的时间是一个区间，
# 顺着时刻算，在区间起点加上这条鱼的贡献，
# 在区间终点减去这条鱼的贡献，求贡献的最大值。
# https://zhuanlan.zhihu.com/p/576489227

# !events扫描线 注意先进后出

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, A = map(int, input().split())
    fish = []
    for _ in range(n):
        weight, pos, speed = map(int, input().split())
        fish.append((weight, pos, speed))

    res = 0
    # 枚举每一条鱼在网的左侧，以这条鱼为参照物
    for i in range(n):
        _, pos1, speed1 = fish[i]
        events = []
        for j in range(n):
            weight2, pos2, speed2 = fish[j]
            posDiff, speedDiff = pos2 - pos1, speed2 - speed1
            if speedDiff == 0:
                if 0 <= posDiff <= A:
                    events.append((0, 0, weight2))
            elif speedDiff < 0:
                left, right = (A - posDiff) / speedDiff, -posDiff / speedDiff
                events.append((left, 0, weight2))
                events.append((right, 1, -weight2))
            else:
                left, right = -posDiff / speedDiff, (A - posDiff) / speedDiff
                events.append((left, 0, weight2))
                events.append((right, 1, -weight2))

        events.sort()
        curSum = 0
        for time, kind, weight in events:
            curSum += weight
            if time >= 0:
                res = max(res, curSum)

    print(res)
