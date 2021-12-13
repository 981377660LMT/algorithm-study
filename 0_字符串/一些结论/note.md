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
######################################
# Wiggle sort
nums = sorted(nums)
half = (len(nums) + 1) >> 1
nums[::2], nums[1::2] = nums[:half], nums[half:]
return nums
```

3. 分成两半 前半长度`不大于`后半长度

```Python
s='asasas'
half=len(s)//2

left=s[:half]
right=s[half:]
```

4. 枚举插入字符串每个位置的字符

```Python
for addIndex in range(len(w)):
    nextChar=w[:addIndex] + addChar + w[addIndex:]
```

5. 枚举分割非空子串

```Python
for i in range(1,len(s)):
    if backtrack(s[i:], int(s[:i])):
        return True
```
