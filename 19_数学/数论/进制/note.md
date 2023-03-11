任意进制字符串相互转换的通用解法(注意特判 0)
`母题-数的进制转换`

1. 十进制转其他进制 toString 函数 ：正序 divmod
   `405. 数字转换为十六进制数`

```JS

 const res = []

 while (num) {
    const digit = num % radix
    res.push(charByDigit[digit])
    num = ~~(num / radix)
  }

 return res.reverse().join('')
```

```Python
def toString(num: int, radix: int) -> str:
    """将数字转换为指定进制的字符串"""
    if num < 0:
        return '-' + toString(-num, radix)

    if num == 0:
        return '0'

    res = []
    while num:
        div, mod = divmod(num, radix)
        res.append(charByDigit[mod])
        num = div
    return ''.join(res)[::-1] or '0'
```

2. 其他进制转十进制 parseInt 函数 ：倒序相加

`171. Excel 表列序号-其他进制转十进制`

```Python
res = 0
base = 1
for i in range(len(string) - 1, -1, -1):
   char = string[i]
   res += base * digitByChar[char]
   base *= RADIX
return res
```

---

结论
