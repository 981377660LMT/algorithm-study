# https://atcoder.jp/contests/abc171/tasks/abc171_f

# 输入 k(≤1e6) 和一个长度不超过 1e6 的字符串 s，由小写字母组成。
# !你需要在 s 中插入恰好 k 个小写字母。
# 输出你能得到的字符串的个数，模 1e9+7。


# https://blog.csdn.net/cqbzlydd/article/details/118608222
# https://www.acwing.com/solution/content/147772/
# !隔板法
# !本题可以等价转化为，有多少长度为 len(s) + k  的字符串存在子序列s 。
# 为了避免算重，我们规定如果有相同字符，则取序列中的最后一个 。
# 例如：字符串 abbcddc 的子序列 bc，从字符串的结尾往开头找，
# c 找到字符串末尾，b 找到字符串第三位。
# 不难发现，对于一个给定的字符串所对应的唯一子序列是确定的。
# !1.枚举开头位置i,前面的数没有限制，一共有 26^i 种情况。
# !2.对于后面n+k-i-1个位置，需要选出n-1个数确定子序列在字符串中的位置(作为最后一个)
# !3.对于剩下k-i个不在子序列中的数,不能与前一个子序列s中的字符一样相同，一共有 25^(k-i) 种情况。

MOD = int(1e9 + 7)
N = int(2e6 + 10)
fac = [1] * N
ifac = [1] * N
for i in range(1, N):
    fac[i] = (fac[i - 1] * i) % MOD
    ifac[i] = (ifac[i - 1] * pow(i, MOD - 2, MOD)) % MOD


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def strivore(s: str, k: int) -> int:
    n = len(s)
    res = 0
    for i in range(k + 1):
        count1 = pow(26, i, MOD)
        count2 = C(n + k - i - 1, n - 1)
        count3 = pow(25, k - i, MOD)
        res = (res + ((count1 * count2) % MOD) * count3) % MOD
    return res


if __name__ == "__main__":
    k = int(input())
    s = input()
    print(strivore(s, k))
