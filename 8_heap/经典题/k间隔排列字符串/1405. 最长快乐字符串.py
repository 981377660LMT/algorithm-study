from heapq import heapify, heappop, heappush

# 如果字符串中不含有任何 'aaa'，'bbb' 或 'ccc' 这样的字符串作为子串，那么该字符串就是一个「快乐字符串」。
# s 中 最多 有a 个字母 'a'、b 个字母 'b'、c 个字母 'c' 。
# s 是一个尽可能长的快乐字符串。


# 每次都选择剩下数量最多的字母添加
class Solution:
    def longestDiverseString(self, a: int, b: int, c: int) -> str:
        dic = dict(a=a, b=b, c=c)
        pq = [(-count, char) for char, count in dic.items() if count > 0]
        heapify(pq)
        res = []

        while pq:
            count, char = heappop(pq)
            count = -count

            if len(res) >= 2 and res[-2] == res[-1] == char:
                if pq:
                    next_count, next_char = heappop(pq)
                    next_count = -next_count
                    res.append(next_char)
                    if next_count - 1 > 0:
                        heappush(pq, (1 - next_count, next_char))
                else:
                    break

            res.append(char)
            if count - 1 > 0:
                heappush(pq, (1 - count, char))

        return "".join(res)


print(Solution().longestDiverseString(a=1, b=1, c=7))
# 输出："ccaccbcc"
# 解释："ccbccacc" 也是一种正确答案。
