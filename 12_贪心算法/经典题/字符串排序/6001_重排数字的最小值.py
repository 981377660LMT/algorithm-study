class Solution:
    def smallestNumber(self, num: int) -> int:
        if num == 0:
            return 0
        digits = list(str(abs(num)))

        if num < 0:
            digits.sort(reverse=True)
            return -int(''.join(digits))
        else:
            digits.sort()
            noneZero = next(i for i, digit in enumerate(digits) if digit != '0')
            digits = digits[noneZero : noneZero + 1] + digits[:noneZero] + digits[noneZero + 1 :]
            return int(''.join(digits))

