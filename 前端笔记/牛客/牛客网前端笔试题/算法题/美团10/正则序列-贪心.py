# 有一天小团得到了一个长度为n的任意序列s，他需要在有限次操作内，将这个序列变成一个正则序列，每次操作他可以任选序列中的一个数字，并将该数字加一或者减一。
# 请问他最少用多少次操作可以把这个序列变成正则序列(1到n的排列)？


# 改动最少的方案一定是对输入序列和正则序列中相同排名的元素


n = int(input())
nums = [int(i) for i in input().split()]
nums.sort()
res = 0
for i in range(n):
    res += abs(i + 1 - nums[i])
print(res)
