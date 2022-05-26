import string


allChar = string.digits + string.ascii_lowercase
charByDigit = {i: char for i, char in enumerate(allChar)}


class Solution(object):
    def toHex(self, num):
        """
        :type num: int
        :rtype: str
        """
        if num < 0:
            num += 2 ** 32

        if num == 0:
            return '0'

        res = []
        while num:
            div, mod = divmod(num, 16)
            res.append(charByDigit[mod])
            num = div
        return ''.join(res)[::-1] or '0'
