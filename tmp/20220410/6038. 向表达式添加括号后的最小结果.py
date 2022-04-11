from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)

# 6038. 向表达式添加括号后的最小结果


class Solution:
    def minimizeResult(self, exp: str) -> str:
        n = len(exp)
        pIndex = exp.find('+')
        resMin = int(1e20)
        resStr = ''

        for leftI in range(pIndex):
            for rightI in range(pIndex + 2, n + 1):
                sa, sb, sc, sd = (
                    exp[:leftI],
                    exp[leftI:pIndex],
                    exp[pIndex:rightI],
                    exp[rightI:],
                )

                a = sa or 1
                b = sb
                c = sc
                d = sd or 1

                cur = min(resMin, int(a) * (int(b) + int(c)) * int(d))
                if cur < resMin:
                    resMin = cur
                    resStr = str(sa) + '(' + str(sb) + str(sc) + ')' + str(sd)

        return resStr


print(Solution().minimizeResult("247+38"))

