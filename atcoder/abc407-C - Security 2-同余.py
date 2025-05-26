# C - Security 2
# https://atcoder.jp/contests/abc407/tasks/abc407_c
# 记 S = s₁s₂…sₙ（每 sᵢ 都是 0–9 的字符），我们从空串 t 出发，
# 要通过若干次按 A（末尾加 ‘0’）或按 B（串中所有数字 +1 mod 10）得到 t = S。


def min_presses_to_S(S: str) -> int:
    n = len(S)
    nums = list(map(int, S))
    totalB = 0
    for i in range(n - 1):
        totalB += (nums[i] - nums[i + 1]) % 10
    totalB += nums[-1] % 10
    return n + totalB


if __name__ == "__main__":
    S = input().strip()
    print(min_presses_to_S(S))
