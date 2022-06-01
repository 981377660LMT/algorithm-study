给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。

```Python
# Python
def repeatedSubstringPattern(s: str) -> bool:
    """给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。"""
    period = (s + s).find(s, 1, -1)
    return period != -1
```
