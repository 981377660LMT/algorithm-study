1. 枚举删去字符串每个位置的字符

```Python
for removeIndex in range(len(w)):
    nextChar=w[:removeIndex] + w[(removeIndex + 1) :]
```

2. 分成两半 前半长度不小于后半长度

```Python
s='asasas'
half=len(s+1)//2

left=s[:half]
right=s[half:]
```

3. 枚举插入字符串每个位置的字符

```Python
for addIndex in range(len(w)):
    nextChar=w[:addIndex] + addChar + w[addIndex:]
```
