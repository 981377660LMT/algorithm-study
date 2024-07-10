https://juejin.cn/user/650530414137534
https://zlxiang.com/

---

# [React 使用 hook 判断组件是否卸载](https://juejin.cn/post/6879560507838169096)

组件卸载后，还调用了 setState，造成了内存泄漏
Warning: Can't perform a React state update on an unmounted component. This is a no-op, but it indicates a memory leak in your application. To fix, cancel all subscriptions and asynchronous tasks in %s.%s a useEffect cleanup function.

```js
export const useMounted = () => {
  const mountedRef = useRef(false)
  useEffect(() => {
    mountedRef.current = true
    return () => {
      mountedRef.current = false
    }
  }, [])
  return () => mountedRef.current
}
```
