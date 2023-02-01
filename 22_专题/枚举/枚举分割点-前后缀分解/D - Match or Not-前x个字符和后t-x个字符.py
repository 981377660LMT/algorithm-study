# 给定两个由小写字母和?组成的字符串S,T，
# 对于 x∈[0,|T|]，问S′和 T是否匹配。
# 其中 S′由 S串的前 x个字符和后 |T|−x个字符组成。
# 两个字符串匹配，当且仅当将其中的 ?替换成英文字母后，两个字符串相同。
# len(s)>=len(t)
if __name__ == "__main__":
    s = input()
    t = input()
    m = len(t)

    def isMatch(c1: str, c2: str) -> bool:
        return c1 == "?" or c2 == "?" or c1 == c2

    preSum, sufSum = [0] * (m + 1), [0] * (m + 1)  # 前i个和后i个的diff数量
    for i in range(m):
        preSum[i + 1] = preSum[i] + (not isMatch(s[i], t[i]))
        sufSum[i + 1] = sufSum[i] + (not isMatch(s[~i], t[~i]))

    for x in range(m + 1):
        diff = preSum[x] + sufSum[m - x]  # 前x个和后len(t)-x个的diff数量
        print("Yes" if diff == 0 else "No")
