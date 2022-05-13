1. 看题目 pre 应该是哪个点 就从哪个点开始出发
2. 一般会有

   ```JS
    root.left = ...
    root.right = pre
    pre.right= ...
    pre = root
   ```
