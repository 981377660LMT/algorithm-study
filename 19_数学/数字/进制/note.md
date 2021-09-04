对于一般性的进制转换题目，只需要不断地对 columnNumbercolumnNumber 进行 % 运算取得最后一位，
然后对 columnNumbercolumnNumber 进行 / 运算，将已经取得的位数去掉，
直到 columnNumbercolumnNumber 为 00 即可。

```JS
转其他进制
 const res = []

 while (num) {
    const digit = num % 16
    res.push(arr[digit])
    num = ~~(num / 16)
  }

 return res.reverse().join('')
```
