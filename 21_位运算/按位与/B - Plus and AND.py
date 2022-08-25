# 长为n的非负整数数组
# 至多可以进行M次操作：每次选中一个整数加上1
# 最后选出k个数 最大化这k个数的按位与的值

# !贪心 (2^i>2^0+2^1+2^2+...+2^(i-1)) + 二分
# 从高位开始判断，判断每个数字当前位如果置为1需要多少步，
# 如果当前位原本就是1，则不消耗，如果原本不是，
# 则消耗低位后，需要将低位全部置0
# 然后排序，选消耗最少的k个，如果满足其消耗小于m，
# 直接默认该为置为1，并保留之前的修改和消耗，继续判断下一位

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m, k = map(int, input().split())
nums = list(map(int, input().split()))


def check(mid: int) -> bool:
    cost = [0] * n  # 每个数与mid与运算为mid的改变代价
    for i in range(n):
        num = nums[i]
        for bit in range(32, -1, -1):
            # !需要从这一位发生改变
            if (mid & (1 << bit)) and ((num & (1 << bit)) == 0):
                allOne = (1 << (bit + 1)) - 1
                cost[i] = (mid & allOne) - (num & allOne)
                break

    cost.sort()
    return sum(cost[:k]) <= m


left, right = 1, 1 << 32  # !注意二分边界为left = 1
while left <= right:
    mid = (left + right) // 2
    if check(mid):
        left = mid + 1
    else:
        right = mid - 1

print(right)  # !最右二分
