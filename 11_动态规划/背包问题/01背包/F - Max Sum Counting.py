# 求1-n的子集个数 使得
# !nums1中的最大值不小于nums2的和
# n<=5000
# nums[i]<=5000

# !手がつかない場合は問題の弱点を探るのが良い。
# !排序+枚举最大值
# dp[index][sum] 表示前index个物品中选择物品,重量为sum的方案数
# !01背包:选还是不选
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    pairs = sorted(zip(nums1, nums2), key=lambda x: x[0])

    max_ = max(nums1)
    dp = [0] * (max_ + 1)
    dp[0] = 1
    res = 0
    for a, b in pairs:
        ndp = dp[:]
        for weight in range(max_ + 1):
            pre = weight - b
            if pre >= 0:
                ndp[weight] = (ndp[weight] + dp[pre]) % MOD
            # !选择当前物品时才能更新方案数
            if weight + b <= a:
                res = (res + dp[weight]) % MOD
        dp = ndp

    print(res)
