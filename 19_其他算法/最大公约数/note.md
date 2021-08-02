GCD(Greatest common divisor)
辗转相除法:
GCD(a,b)=GCD(a%c,b)

```Python
def GCD(a: int, b: int) -> int:
    return a if b == 0 else GCD(b, a % b)
时间复杂度：$O(log(max(a, b)))$
```

更相减损术:
GCD(a,b)=GCD(a-b,b) (a>b)

```Python
def GCD(a: int, b: int) -> int:
    if a == b:
        return a
    if a < b:
        return GCD(b - a, a)
    return GCD(a - b, b)
最坏时间复杂度为 O(max(a, b)))
```

a\*b/最小公倍数

```Python
def lcm(x, y):
    return x * y // gcd(x, y)
```
