# 字符串重新排列不含前导0的数字的个数


from math import factorial


def reArrangeDigits(s: str) -> int:
    n = len(s)
    counter = [0] * 10
    for c in s:
        counter[int(c)] += 1

    res = (n - counter[0]) * factorial(n - 1)
    for i in range(10):
        res //= factorial(counter[i])
    return res


if __name__ == "__main__":
    print(reArrangeDigits("1234567890"))  # 362880
    print(reArrangeDigits("123456789"))  # 362880
