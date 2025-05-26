# abc407-A - Approximation-round
# https://atcoder.jp/contests/abc407/editorial/13102
# 输出离a/b最近的整数
#
# 答案就是 round(a/b)
#
# !也可以不用浮点数计算.
# 离a/b最近的整数，就是不超过 a/b + 0.5 的最大整数
# 即 (2a+b)//2b
# 也可变形为 (a + b // 2) // b
if __name__ == "__main__":
    a, b = map(int, input().split())
    print((a + b // 2) // b)
