# https://nyaannyaan.github.io/library/misc/base64.hpp
# base64的转换/base64的加密与解密
# 用于压缩源码(一般是int数组)/压缩打表数据
# !通常配合`部分打表`技巧使用(就是前半部分是打表,后半部分是计算)
# https://nyaannyaan.github.io/library/verify/verify-yuki/yuki-0502-base64.test.cpp


from typing import List


class Base64:
    _BASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

    @staticmethod
    def encode(nums: List[int]) -> str:
        x, y = nums[0], nums[0]
        for z in nums:
            x = max(x, z)
            y = min(y, z)
        n = len(nums)
        b = max(6, 64 if y < 0 else x.bit_length())
        sb = [0] * ((b * n + 11) // 6)
        sb[0] = b
        for i in range(n):
            for j in range(b):
                if nums[i] >> j & 1:
                    sb[(i * b + j) // 6 + 1] |= 1 << ((i * b + j) % 6)
        return "".join(Base64._BASE[c] for c in sb)

    @staticmethod
    def decode(s: str) -> List[int]:
        sb = [Base64._ibase(c) for c in s]
        b = sb[0]
        m = len(sb) - 1
        nums = [0] * (6 * m // b)
        for i in range(m):
            for j in range(6):
                if sb[i + 1] >> j & 1:
                    nums[(i * 6 + j) // b] |= 1 << ((i * 6 + j) % b)
        return nums

    @staticmethod
    def _ibase(c: str) -> int:
        if c >= "a":
            return 0x1A + ord(c) - ord("a")
        if c >= "A":
            return 0x00 + ord(c) - ord("A")
        if c >= "0":
            return 0x34 + ord(c) - ord("0")
        if c == "+":
            return 0x3E
        if c == "/":
            return 0x3F
        return 0x40


if __name__ == "__main__":
    nums = list(range(100))
    s = Base64.encode(nums)
    print(s)
    print(Base64.decode(s))
