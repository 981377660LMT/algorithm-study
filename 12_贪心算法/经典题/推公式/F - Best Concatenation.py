# 有n个字符串，请构造一种字符串使得该字符串为n的字符串的任意顺序的串联，
# 最大化其价值。字符串的价值定义为：每有一个数，
# !其贡献为该数之前出现过的X的数量 * 该数字，价值为所有数的贡献和。
# 例如 `XXX1X359`
# (X,1) 有3组 (X,3) (X,5) (X,9) 各有4组
# 总分为1*3+3*4+5*4+9*4=71

# !贪心推公式+排序
# !注意到交换相邻的两个字符只会改变这两个字符的贡献，
# 因此我们可以贪心地交换相邻的两个字符，使得贡献最大化。
# x[a]表示a字符串内x的数量，sum[b]表示b字符串内数字的和，
# x[b]表示b字符串内x的数量。那么假如a排在b前面，
# 字符串之间产生的贡献就是sum[b] * x[a]，
# 加入b排在a前面，贡献为sum[a] * x[b]。
# !也就是说，如果sum[b] * x[a] > sum[a] * x[b]更大，
# 那么a更要排在前面，我们按照这个排序，然后计算所有的贡献即可。
from functools import cmp_to_key
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
words = [(i, input()) for i in range(n)]  # (index, word)


countX = [0] * n
digitSum = [0] * n
for i, (_, word) in enumerate(words):
    for char in word:
        if char == "X":
            countX[i] += 1
        else:
            digitSum[i] += int(char)

words.sort(
    key=cmp_to_key(lambda a, b: digitSum[b[0]] * countX[a[0]] - digitSum[a[0]] * countX[b[0]]),
    reverse=True,
)


preSum = [0]
for _, word in words:
    for char in word:
        preSum.append(preSum[-1] + int(char == "X"))

res, index = 0, 0
for _, word in words:
    for char in word:
        if char != "X":
            res += int(char) * preSum[index]
        index += 1
print(res)
