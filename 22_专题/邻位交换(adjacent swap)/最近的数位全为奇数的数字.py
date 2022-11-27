# 最近的数位全为奇数的数字
# Next Closest Odd Digit Number

# 0 ≤ n ≤ 2 ** 31 - 1
class Solution:
    def solve(self, num):
        string = str(num)
        n = len(string)

        # 找到第一个偶数位
        evenCand = next((i for i in range(n) if int(string[i]) % 2 == 0), -1)
        if evenCand == -1:
            return num

        bigCand = str(int(string[: evenCand + 1]) + 1)
        if string[evenCand] == "0":
            if string[evenCand - 1] == "1":
                smallCand = str(int(string[: evenCand + 1]) - 1)  # 110-1 = 109
            else:
                smallCand = str(int(string[: evenCand + 1]) - 11)  # 130-11= 119
        else:
            smallCand = str(int(string[: evenCand + 1]) - 1)

        # tail "1" and "9"
        for _ in range(evenCand + 1, n):
            bigCand += "1"
            smallCand += "9"

        bigCand, smallCand = int(bigCand), int(smallCand)
        diff1 = num - smallCand
        diff2 = bigCand - num
        if diff1 < diff2:
            return smallCand
        elif diff1 >= diff2:
            return bigCand
