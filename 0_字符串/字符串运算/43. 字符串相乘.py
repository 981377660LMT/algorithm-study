# 字符串相乘 模拟乘法
class Solution:
    def multiply(self, num1: str, num2: str) -> str:
        if num1 == "0" or num2 == "0":
            return "0"

        n1, n2 = len(num1), len(num2)
        res = [0] * (n1 + n2)
        for i1 in range(n1 - 1, -1, -1):
            a = int(num1[i1])
            for i2 in range(n2 - 1, -1, -1):
                b = int(num2[i2])
                res[i1 + i2 + 1] += a * b

        # 处理进位
        for i in range(n1 + n2 - 1, 0, -1):
            res[i - 1] += res[i] // 10
            res[i] %= 10

        if res[0] == 0:
            return "".join(map(str, res[1:]))
        return "".join(map(str, res))
