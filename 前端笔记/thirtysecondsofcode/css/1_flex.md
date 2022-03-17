flex : flex-grow flex-shrink flex-basis

flex-grow 决定了如果需要的话，这个东西能长多大

- 接受单个正数(无单位) ，默认值为 0

flex-shrink 决定了如果有必要的话，这些东西可以缩小多少

- 接受单个正数(无单位) ，默认值为 1

flex-basis 决定**在分配剩余空间之前**，确定 flex 项的初始大小

![图解](https://www.30secondsofcode.org/assets/blog_images/flexbox-diagram.webp)

属性
**display: flex or display: inline-flex**:
当 Flex Box 容器没有设置宽度大小限制时，当 display 指定为 flex 时，FlexBox 的宽度会填充父容器，当 display 指定为 inline-flex 时，FlexBox 的宽度会包裹子 Item，如下图所示：
**flex-direction**：
row/row-reverse
column/column-reverse
**flex-wrap**
Nowrap (默认值) : 所有 flex 项将在一行上
wrap：Flex 项将包装到多行上，从上到下
**justify-content**
flex-start/flex-end/center/space-between/space-around
**align-items**
flex-start/flex-end/center

项目的属性
**align-self**
属性允许单个项目有与其他项目不一样的对齐方式，可覆盖 align-items 属性

---
