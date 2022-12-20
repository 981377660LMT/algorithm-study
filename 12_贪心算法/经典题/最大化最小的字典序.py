# 最大化最小字典序(最大化轮转后字典序最小的字符串的字典序)
# https://atcoder.jp/contests/code-festival-2017-qualb/tasks/code_festival_2017_qualb_f


# iehnさん:
# !辞書順最小を最大化したいだけなんだから、
# !その時点で最小のに最大のをつなげていけばいいだけ
def largestSmallestCyclicShift(x: int, y: int, z: int) -> str:
    # 由x个a + y个b + z个c组成的字符串，轮转后字典序最小的字符串记为s
    # 请最大化s的字典序
    # x+y+z <=50
    sb = ["a"] * x + ["b"] * y + ["c"] * z
    while len(sb) > 1:
        sb.sort()
        sb[0] += sb.pop()
    return sb[0]


x, y, z = map(int, input().split())
print(largestSmallestCyclicShift(x, y, z))
