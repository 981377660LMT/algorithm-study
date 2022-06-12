from itertools import groupby
from string import ascii_lowercase, ascii_uppercase, digits

from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def strongPasswordCheckerII(self, password: str) -> bool:
        n = len(password)
        wSet = set(password)
        if n < 8:
            return False
        groups = [[char, len(list(group))] for char, group in groupby(password)]
        if len(groups) < n:
            return False
        if all(char not in wSet for char in ascii_lowercase):
            return False
        if all(char not in wSet for char in ascii_uppercase):
            return False
        if all(char not in wSet for char in digits):
            return False
        if all(char not in wSet for char in "!@#$%^&*()-+"):
            return False
        return True

