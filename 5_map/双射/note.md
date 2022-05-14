要么同时不存在 要么同时都存在且相等

```Python
def isIsomorphic(s: str, t: str) -> bool:
    if len(s) != len(t):
        return False
    return len(set(s)) == len(set(t)) == len(set(zip(s, t)))
```
