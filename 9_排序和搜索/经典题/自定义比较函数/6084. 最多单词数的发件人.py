from collections import Counter
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def largestWordCount2(self, messages: List[str], senders: List[str]) -> str:
        counter = Counter()
        for msg, sender in zip(messages, senders):
            counter[sender] += len(msg.split())
        # 这一段可以用max+key来写 要求的属性用max 然后key lambda x: x... 来限定属性
        max_ = max(counter.values())
        res = []
        for key in counter:
            if counter[key] == max_:
                res.append(key)
        res.sort(reverse=True)
        return res[0]

    def largestWordCount(self, messages: List[str], senders: List[str]) -> str:
        counter = Counter()
        for msg, sender in zip(messages, senders):
            counter[sender] += len(msg.split())
        # 这一段可以用max+key来写 要求的属性用max 然后key lambda x: x... 来限定属性
        # max的比较器可以是元组
        return max(counter, key=lambda x: (counter[x], x))

