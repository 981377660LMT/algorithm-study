# 166. 分数到小数
# 将分数转换为小数，如果小数部分为循环小数，则将循环的部分括在括号内。


class Solution:
    def fractionToDecimal(self, numerator: int, denominator: int) -> str:
        sb = []
        a, b = numerator, denominator
        if a == 0:
            return "0"
        if a > 0 ^ b > 0:
            sb.append("-")

        a, b = abs(a), abs(b)
        sb.append(str(a // b))
        if a % b == 0:
            return "".join(sb)

        sb.append(".")
        memo = dict()  # 被除数: 位置
        while True:
            a = (a % b) * 10
            if a == 0 or a in memo:
                break
            memo[a] = len(sb)
            sb.append(str(a // b))

        if a == 0:
            return "".join(sb)

        cycleStart = memo[a]
        sb.insert(cycleStart, "(")
        sb.append(")")
        return "".join(sb)


if __name__ == "__main__":
    print(Solution().fractionToDecimal(1, 29))
