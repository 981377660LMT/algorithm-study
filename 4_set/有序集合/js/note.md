`传 undefined 或不传都会触发默认参数`

不会死循环

```JS
function* inorder(root: Node | null | undefined = demo): Generator<number> {
  if (!root) return
  yield* inorder(null)
  yield root.val
  yield* inorder(null)
}
```

会死循环

```JS
function* inorder(root: Node | null | undefined = demo): Generator<number> {
  if (!root) return
  yield* inorder(undefined)
  yield root.val
  yield* inorder(undefined)
}

function* inorder(root: Node | null | undefined = demo): Generator<number> {
  if (!root) return
  yield* inorder()
  yield root.val
  yield* inorder()
}
```

解决这个问题的方法是空值`传 null`避免触发默认参数，或者统一不用默认参数
