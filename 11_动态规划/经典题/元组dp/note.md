1. 注意有的题元组转移时要**排序** 否则无法起到记忆化的作用
   [1655. 分配重复整数](1655.%20%E5%88%86%E9%85%8D%E9%87%8D%E5%A4%8D%E6%95%B4%E6%95%B0.py)
   [1900. 最佳运动员的比拼回合-元组为记忆化参数的 dfs](1900.%20%E6%9C%80%E4%BD%B3%E8%BF%90%E5%8A%A8%E5%91%98%E7%9A%84%E6%AF%94%E6%8B%BC%E5%9B%9E%E5%90%88-%E5%85%83%E7%BB%84%E4%B8%BA%E8%AE%B0%E5%BF%86%E5%8C%96%E5%8F%82%E6%95%B0%E7%9A%84dfs.py)

   **!注意这个 sorted 把排列变成了组合**
   **用 tuple 本质就是子集状压 dp**

   ```Python
   # !注意这里要保持顺序 否则就起不到记忆化的效果
   nextRemain = sorted(remain[:i] + (num - quantity[index],) + remain[i + 1 :])
   dfs(index + 1, tuple(nextRemain))


   # !注意这里要排序
   nexts = dfs(tuple(sorted(winners)))
   ```
