# C - Max GCD 2
# 最大gcd
# 给定上下界,求lower<=x<y<=upper,使得gcd(x,y)最大
# 1<=lower<upper<=2e5

# !枚举答案cand 看是否能满足gcd(x,y)=cand


def maxGcd(lower: int, upper: int) -> int:
    res = 1
    for cand in range(1, int(2e5) + 1):
        count = upper // cand - (lower - 1) // cand  # [lower,upper]内cand的倍数个数
        if count >= 2:
            res = cand
    return res


A, B = map(int, input().split())
print(maxGcd(A, B))
