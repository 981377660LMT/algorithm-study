# 给定一个字符串s，s由‘0’~‘9’组成，设d是s中最大的数字，
# 不同的进位制n(n>=d+1)使得s在该进位制下的值小于等于M,求出不同的值的个数

# !二分最大进制,注意特判s长度为1的情况


def baseN(s: str, upper: int) -> int:
    def check(mid: int) -> bool:
        """将s看成是mid进制的数,是否小于等于upper"""
        cur = 0
        for c in s:
            cur = cur * mid + int(c)
            if cur > upper:
                return False
        return True

    if len(s) == 1:  # !此时进制不同,值的个数相同(1或0个值)
        return 1 if int(s) <= upper else 0

    max_ = max(int(c) for c in s)
    left, right = max_ + 1, int(1e18)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    return right - (max_ + 1) + 1


if __name__ == "__main__":
    s = input()
    M = int(input())
    print(baseN(s, M))
