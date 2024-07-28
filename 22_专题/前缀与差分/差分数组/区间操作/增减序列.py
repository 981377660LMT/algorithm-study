# 100. 增减序列
# https://www.acwing.com/file_system/file/content/whole/index/content/3651/
# https://www.acwing.com/solution/content/5060/
# 不变量/差分数组
# 给定一个长度为 n 的数列 a1,a2,…,an，每次可以选择一个区间 [l,r]，使下标在这个区间内的数都加一或者都减一。
# 求至少需要多少次操作才能使数列中的所有数都一样，并求出在保证最少次数的前提下，最终得到的数列可能有多少种。
#
# 选取一个区间[i, j] diff[i] ++ diff[j + 1] -- 或 diff[i] -- diff[j + 1] ++ 两操作异号
# 使diff[2] ~ diff[n] 全变成0，所有数就和diff[1]一样了
# !操作的时候，要找两个数配对，那么 负数++，正数--，是不是就最快了


n = int(input())
nums = [int(input()) for _ in range(n)]


diff = [0] + [b - a for a, b in zip(nums, nums[1:])]
pos = 0
neg = 0
for d in diff:
    if d > 0:
        pos += d
    elif d < 0:
        neg += -d

print(max(pos, neg))  # 最少操作次数
print(abs(pos - neg) + 1)  # 最终得到的数列
