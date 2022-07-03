![动态数组1e7正好超出](image/js动态数组长度上限9e6/1653299288430.png)

![动态数组最多存9e6个数](image/js动态数组长度上限9e6/1653299342826.png)

![9e6动态数组扩容耗时约80ms](image/js动态数组长度上限9e6/1653299450441.png)

![静态数组没有扩容时间花费](image/js动态数组长度上限9e6/1653299969803.png)

![Uint32Array最多存3.5e7个数,140MB](image/js动态数组长度上限9e6/1653300105875.png)
`print(3.5 * int(1e7) * 32/8 // int(1e6))` => **140MB**
![Uint8Array最多存14e7个数](image/js动态数组长度上限9e6/1653300233176.png)

```JS
/**
 * @param {number[][]} forest
 * @return {number}
 */
var cutOffTree = function (forest) {
  const list = new Uint8Array(3.5e7)  // 静态数组
  let x = 0
  for (let i = 0; i < 1e9; i++) {
    list[i] = 0
    x += i * (i + 1)
  }
  return 0
}
```

**可以看到 动态数组耗费的空间是 Uint32Array 的 4 倍左右**
