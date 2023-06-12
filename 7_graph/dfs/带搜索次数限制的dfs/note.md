带搜索次数限制的 dfs

```python
# dfs内部多了一行剪枝的逻辑
nonlocal count
if count >= k - 1:  # !限制dfs次数
    return
```
