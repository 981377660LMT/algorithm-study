import string


allChar = string.digits + string.ascii_lowercase
charByDigit = {i: char for i, char in enumerate(allChar)}


class Solution:
    def convertToBase7(self, num: int) -> str:
        if num < 0:
            return '-' + self.convertToBase7(-num)

        if num == 0:
            return '0'

        res = []
        while num:
            div, mod = divmod(num, 7)
            res.append(charByDigit[mod])
            num = div
        return ''.join(res)[::-1] or '0'
