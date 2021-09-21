离线排序 vs 在线排序
离线排序:
实现将 query 按照严格到宽松的顺序排序，逐个建图

常用预处理：

```JS
queries = queries.map((v, i) => [...v, i]).sort((a, b) => a[2] - b[2])
```
