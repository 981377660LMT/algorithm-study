# https://atcoder.jp/contests/abc293/tasks/abc293_f
# 给定一个正数2<=n<=10^18
# !问有多少个大于等于2的进制radix使得在这个进制下,n的每位数都是0或1

# 解:
# !先把可能的进制算出来 再逐个检查
# !若某个i+1(i>=2)位数上有解，那么B=floor(N^(1/i)) 是唯一可能的解,即 B^i<=N<B^(i+1)
# 如果B更大，就N<B^i，达不到这个位数；
# 如果B更小，用二项式定理式子可以推出系数显然不全为1


def check(num: int, radix: int) -> bool:
    """将数字转换为指定进制的字符串后是否每位都是0或1"""
    if radix <= 1:
        return False
    while num:
        div, mod = num // radix, num % radix
        if mod > 1:
            return False
        num = div
    return True


def zeroOrOne(n: int) -> int:
    toCheck = set([n, n - 1])  # 10 11
    for i in range(2, n.bit_length() + 1):
        x = int(n ** (1 / i))
        toCheck.add(x)
        toCheck.add(x + 1)
    return sum(check(n, radix) for radix in toCheck)


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n = int(input())
        print(zeroOrOne(n))
