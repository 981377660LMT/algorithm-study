# 给你一个字符串 s 。s[i] 要么是小写英文字母，要么是问号 '?' 。
# 对于长度为 m 且 只 含有小写英文字母的字符串 t .
# !我们定义函数 cost(i) 为下标 i 之前（也就是范围 [0, i - 1] 中）出现过与 t[i] 相同 字符出现的次数。
# 你的任务是用小写英文字母 替换 s 中 所有 问号，使 s 的 分数最小 。
# 请你返回替换所有问号 '?' 之后且分数最小的字符串。
# 如果有多个字符串的 分数最小 ，那么返回字典序最小的一个。


# 根据cost，每个字母地贡献是(freq - 1) * freq // 2
# 为了使得cost最小，我们需要尽量使得freq接近
# 考虑用一个最小堆来维护freq

from collections import Counter
from heapq import heapify, heappop, heappush
from string import ascii_lowercase


class Solution:
    def minimizeStringValue(self, s: str) -> str:
        freq = Counter(s)
        pq = [(freq[c], c) for c in ascii_lowercase]
        heapify(pq)

        toadd = []
        for _ in range(s.count("?")):
            f, c = heappop(pq)
            toadd.append(c)
            heappush(pq, (f + 1, c))
        toadd.sort()  # 字典序最小

        sb = list(s)
        ptr = 0
        for i in range(len(sb)):
            if sb[i] == "?":
                sb[i] = toadd[ptr]
                ptr += 1
        return "".join(sb)
