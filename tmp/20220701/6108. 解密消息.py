import string

# 注意要 iter 把 iterable 套一下


class Solution:
    def decodeMessage(self, key: str, message: str) -> str:
        """迭代器简化了'取出'下一个元素的逻辑

        注意iterable变为迭代器 要用iter包装
        """
        mapping, it = dict(), iter(string.ascii_lowercase)
        for char in key:
            if char not in mapping and char != " ":
                mapping[char] = next(it)
        return "".join(mapping.get(c, c) for c in message)


print(
    Solution().decodeMessage(
        key="the quick brown fox jumps over the lazy dog", message="vkbs bs t suepuv"
    )
)
