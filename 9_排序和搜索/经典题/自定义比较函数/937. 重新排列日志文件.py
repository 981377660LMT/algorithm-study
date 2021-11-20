# 字母日志：除标识符之外，所有字均由小写字母组成
# 数字日志：除标识符之外，所有字均由数字组成

#  请按下述规则将日志重新排序：

#  所有 字母日志 都排在 数字日志 之前。
#  字母日志 在内容不同时，忽略标识符后，按内容字母顺序排序；在内容相同时，按标识符排序。
#  数字日志 应该保留原来的相对顺序。


# reorderLogFiles(['dig1 8 1 5 1', 'let1 art can', 'dig2 3 6', 'let2 own kit dig', 'let3 art zero'])

#  输出：["let1 art can","let3 art zero","let2 own kit dig","dig1 8 1 5 1","dig2 3 6"]
#  解释：
#  字母日志的内容都不同，所以顺序为 "art can", "art zero", "own kit dig" 。
#  数字日志保留原来的相对顺序 "dig1 8 1 5 1", "dig2 3 6" 。
from typing import Any, List, Tuple


class Solution:
    def reorderLogFiles(self, logs: List[str]) -> List[str]:
        def compare_key(log: str) -> Tuple[Any, ...]:
            left, right = log.split(" ", 1)
            if right[0].isalpha():
                return (0, right, left)
            else:
                return (1,)

        return sorted(logs, key=compare_key)

