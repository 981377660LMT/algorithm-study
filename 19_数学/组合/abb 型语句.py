# leafee 最近爱上了 abb 型语句，比如“叠词词”、“恶心心”
# leafee 拿到了一个只含有小写字母的字符串，她想知道有多少个 "abb" 型的子序列？
# 1≤n≤10^5

# 对每个位置的字符，计算以他开头的贡献
from math import comb
from collections import Counter

n = int(input())
string = input()

counter = Counter(string)
res = 0
for char in string:
    counter[char] -= 1
    for pair in counter.keys():
        if char != pair:
            res += comb(counter[pair], 2)
print(res)

