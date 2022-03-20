<!-- 关键是把下一个中间件放在next参数的位置 -->

1. express 采用 reduce 串联中间件

```JS
return funcs.reduce((a, b) => (req, res, next) => a(req, res, () => b(req, res, next)))
```

2. koa 采用 dfs 逐层调用中间件

```JS
async function dfs(index: number): Promise<void> {
  if (index == middlwwares.length) return
  const middleware = middlwwares[index]
  try {
    await middleware(context, () => dfs(index + 1))
  } catch (error) {
    throw error
  }
}
```
