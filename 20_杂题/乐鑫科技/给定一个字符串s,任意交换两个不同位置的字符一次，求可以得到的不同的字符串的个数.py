# 给定一个字符串s,任意交换两个不同位置的字符一次，求可以得到的不同的字符串的个数
# 1<=len(s)<=1e5
from collections import Counter

string = input()
counter = Counter(string)
n = len(string)

res = 0
canKeep = False

for count in counter.values():
    # 每个组都可以和其他字符交换，每个对算了两次
    res += count * (n - count)
    if count > 1:
        canKeep = True  # 最开始的那个算不算

print(res // 2 + int(canKeep))

