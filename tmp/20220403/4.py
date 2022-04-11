from collections import defaultdict
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 至多调用 encrypt 和 decrypt 总计 200 次


class Encrypter:
    def __init__(self, keys: List[str], values: List[str], dictionary: List[str]):
        self.keyToValue = defaultdict(str)
        for k, v in zip(keys, values):
            self.keyToValue[k] = v

        self.valueToKey = defaultdict(set)
        for value, key in zip(values, keys):
            self.valueToKey[value].add(key)

        self.ok = set(dictionary)

    def encrypt(self, word1: str) -> str:
        """1 <= word1.length <= 2000"""

        res = []
        for c in word1:
            res.append(self.keyToValue[c])
        return ''.join(res)

    def decrypt(self, word2: str) -> int:
        """1 <= word2.length <= 200  是偶数"""

        def check(target) -> bool:
            for i, c in enumerate(target):
                if i >= len(res) or c not in res[i]:
                    return False
            return True

        res = []
        for i in range(0, len(word2), 2):
            cur = word2[i : i + 2]
            res.append(self.valueToKey[cur])

        return sum(check(target) for target in self.ok if len(target) == len(res))


# @lru_cache(None)
# def dfs(index: int) -> bool:
#     ...


ec = Encrypter(["a"], ["pq"], ["a", "ax"])


print(ec.decrypt("pq"))

# for str in [["aaaa"], ["aa"], ["aaaa"], ["aaaaaaaa"], ["aaaaaaaaaaaaaa"], ["aefagafvabfgshdthn"]]:
#     print(ec.decrypt(str[0]))

# 预期：
# [null,1,0,1,1,1,0]

