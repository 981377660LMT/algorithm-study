# 正整数 N が与えられます。
# A×B+C=N を満たす正整数の組 (A,B,C) はいくつありますか？
# 2<=n<=1e6
# !枚举A,可以求出B的范围,从而C也随之确定 O（n）

n = int(input())
res = 0
for a in range(1, n):
    res += (n - 1) // a  # b的范围是[1, (n-1)//a]
print(res)
