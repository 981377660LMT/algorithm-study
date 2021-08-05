1. 快慢指针（两个指针步长不同）

```Python
l = 0
r = 0
while 没有遍历完
  if 一定条件
    l += 1
  r += 1
return 合适的值
```

2. 左右端点指针（两个指针分别指向头尾，并往中间移动，步长不确定）

```Python
l = 0
r = n - 1
while l < r
  if 找到了
    return 找到的值
  if 一定条件1
    l += 1
  else if  一定条件2
    r -= 1
return 没找到
```

3. 固定间距指针（两个指针间距相同，步长相同）

```Python
l = 0
r = k
while 没有遍历完
  自定义逻辑
  l += 1
  r += 1
return 合适的值
```
