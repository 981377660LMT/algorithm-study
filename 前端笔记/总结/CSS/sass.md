1. Sass 嵌套（Nested）

```scss
.father {
  color: red;
  .child {
    color: green;
    &:hover {
      color: red;
    }
    &:active {
      color: blue;
    }
    &-item {
      color: orange;
    }
  }
}
```

嵌套规则很有用很方便，但是你很难想象它实际会生成多少 CSS 语句，嵌套的越深，那么编译为 CSS 的语句就越多，同时消耗的资源也会越多，所以开发者尽量不要嵌套特别深的层级！
