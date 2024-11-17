# abc379-E - Sum of All Substrings-所有子串的数值之和(进制转换)
# https://atcoder.jp/contests/abc379/tasks/abc379_e
#
# 给定一个长达1e5的数字字符串s。定义f(i,j)=s[i..j]。即子串的数字。
# 求∑ni=1∑nj=1f(i,j)。
#
# 贡献法.
# 1. 考虑每个s[i]在答案的贡献。但是这样数字很大，不好处理。
# !2. 考虑最终答案的每一位的数字是多少，即因子1e0,1e1,1e2的系数分别是多少。
# !最终的答案为 ∑10^(n-1-i)*A[i]，其中 A[i]= ∑(j+1)*s[j] (j<=i)。


def sumOfAllSubstrings(n: int, s: str) -> str:
    nums = [int(c) for c in s]
    digits = []
    sum_ = 0
    for i, v in enumerate(nums):
        sum_ += (i + 1) * v
        digits.append(sum_)
    carry = 0
    for i in range(n - 1, -1, -1):
        digits[i] += carry
        carry = digits[i] // 10
        digits[i] %= 10
    if carry:
        digits.insert(0, carry)
    return "".join(map(str, digits))


if __name__ == "__main__":
    N = int(input())
    S = input()
    print(sumOfAllSubstrings(N, S))
