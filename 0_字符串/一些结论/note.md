1. 枚举删去字符串每个位置的字符

```Python
for removeIndex in range(len(w)):
    nextChar=w[:removeIndex] + w[(removeIndex + 1) :]
```

2. 分成两半 前半长度`不小于`后半长度

```Python
s='asasas'
half=len(s+1)//2

left=s[:half]
right=s[half:]
```

2. 分成两半 前半长度`不大于`后半长度

```Python
s='asasas'
half=len(s)//2

left=s[:half]
right=s[half:]
```

3. 枚举插入字符串每个位置的字符

```Python
for addIndex in range(len(w)):
    nextChar=w[:addIndex] + addChar + w[addIndex:]
```

4. 枚举分割非空子串

```Python
for i in range(1,len(s)):
    if backtrack(s[i:], int(s[:i])):
        return True
```
