import string


class Solution:
    def checkInclusion(self, s1: str, s2: str) -> bool:
        if len(s1) > len(s2):
            return False
        target_counter = {k: 0 for k in string.ascii_lowercase}
        cur_counter = {k: 0 for k in string.ascii_lowercase}
        for char in s1:
            target_counter[char] += 1
        for i in range(len(s2)):
            cur_counter[s2[i]] += 1
            if i >= len(s1):
                cur_counter[s2[i - len(s1)]] -= 1
            if cur_counter == target_counter:
                return True
        return False

