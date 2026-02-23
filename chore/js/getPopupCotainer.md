https://github.com/ant-design/ant-design/issues/20367
https://github.com/ant-design/ant-design/issues/49745

---

Ant Design 的 `Tour`（漫游式引导）或其他基于 Popper.js / `@rc-component/trigger` 的浮层组件（如 Tooltip, Popover），在父容器使用了 `transform: scale(...)` 进行缩放时，经常会出现定位偏移（错位）。

这个问题非常经典，以下是其本质原因及基于场景的最佳修复实践。

### 1. 本质原因深度解析

#### 坐标系的不一致

浮层组件（如 Antd 的 Tour）通常渲染在 `document.body` 下（通过 React Portal），这意味着它是基于 **视口（Viewport）** 或 **文档流** 进行绝对定位的。

而目标元素（Target）位于一个被 `transform: scale(0.5)`（举例）影响的容器内。

1.  **getBoundingClientRect 的测量值：**
    当你调用目标元素的 `getBoundingClientRect()` 时，浏览器返回的是它在当前视口中的**视觉坐标**和尺寸。如果容器缩小了 50%，这个矩形的 `width`, `height`, `top`, `left` 都会反映这个缩放后的结果。

2.  **定位逻辑的各异：**
    - **常规逻辑：** 浮层组件通常计算 `目标元素.left + window.scrollX` 来决定自己的 `left`。
    - **Scale 的破坏性：** 虽然 `getBoundingClientRect` 返回了“正确”的视觉位置，但浮层自身通常是一个不受 `scale` 影响的 DOM 节点（因为它挂载在 body 上）。
    - **位移偏差：** 很多定位库（早期的 Popper.js 或某些简易实现）假定 1px 的位移等于 1px 的视觉距离。但在 scale 上下文中，父级容器产生的滚动距离（Scroll）或者 CSS `translate` 的位移，在视觉上被缩放了，但定位库在计算绝对位置时可能没有除以这个 `scale` 系数，导致计算出的 `top/left` 坐标与实际视觉坐标不匹配。

#### 复合层（Compositing Layer）与 Stack Context

`transform` 属性会创建新的层叠上下文（Stacking Context）。如果浮层没有渲染到 body，而是渲染在这个 scale 容器内部，那么浮层也会被缩放，且可能受到父级 `overflow: hidden` 的影响被裁剪。因此，Antd 默认将 Tour 渲染到 body 是有道理的，但这引发了上述的坐标系分离问题。

---

### 2. 修复的最佳实践

根据不同的场景，有以下三种主要的修复策略：

#### 方案一：不仅仅缩放容器，而是使用 `zoom` (不推荐但往往有效)

如果你的场景允许，使用 CSS `zoom` 属性代替 `transform: scale`。`zoom` 会触发布局更新，浏览器会自动处理坐标计算。

- **缺点：** `zoom` 是非标准属性（主要在 Chrome/Edge 有效，Firefox 不支持），且性能不如 transform。

#### 方案二：强制浮层渲染在 Scale 容器内部 (特定场景)

如果将 Tour 渲染在被缩放的容器内部，它就会自动继承 scale 属性，位置通常会由 CSS 布局引擎自动对齐。

可以使用 `getPopupContainer` 属性将浮层挂载节点指定为那个缩放的 div。

```typescript
// 在 Tour 组件上使用
<Tour
  // ...other props
  getPopupContainer={triggerNode => triggerNode.parentElement}
/>
```

- **风险：** 浮层也会被“缩小”或“放大”（文字变小/变大）。如果 scale 是为了适配大屏的大比例缩小，浮层可能会小到看不清。
- **风险：** 父容器必须没有 `overflow: hidden`，否则引导层会被切断。

#### 方案三：在该场景下禁用 Portal 并手动修正 Scale (最稳健的方案)

这是针对复杂 Scale 场景（如数据大屏、低代码编辑器）最通用的解法。我们需要让计算逻辑感知到 Scale 的存在。

如果你使用的是较新版本的 Ant Design (v5+)，底层依赖的 `rc-util/Dom/css.ts` 或浮层定位逻辑可能并未完全覆盖所有 transform 矩阵场景。

**核心思路：禁用 Portal 并反向抵消 Scale**

1.  **挂载在容器内：** 设置 `getPopupContainer` 为缩放容器。
2.  **反向 Scale：** 因为挂载在容器内，Tour 自身也会被缩放。我们需要在 Tour 的样式上应用 `scale(1 / factor)` 来抵消视觉缩放，确保护罩和弹窗大小正常。

但 Antd 的 `Tour` 组件比较特殊，它包含“蒙层（Mask）”和“弹窗（Panel）”。蒙层是 SVG 实现的。

**针对 Antd 5.x Tour 的具体修复代码：**

可以通过 `scale` 能够感知上下文的思路，通常不需要我们手动写样式，而是**确保 Tour 的目标元素计算**包含了 scale。

如果在 Antd Issues 49745 中提到的场景，最有效的“黑科技”修复通常是**劫持 getBoundingClientRect** 或者 **使用 Antd 的 `zIndex` 配合 `getPopupContainer`**。

对于纯粹的错位问题，如果无法修改 Antd 源码，建议使用以下 CSS 配合将组件挂载到非 Scale 区域（或者反过来）。

**推荐：使用 `getPopupContainer` 锁定范围，修正偏移**

如果你的页面结构是：
`Body -> Div(Scale=0.5) -> Button(Target)`

尝试将 Tour 挂载到 Scale 容器层：

```tsx
import React, { useRef } from 'react'
import { Tour, Button } from 'antd'

const App = () => {
  const ref1 = useRef(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // 假设 scaleValue 是当前的缩放比例，例如 0.5
  const scaleValue = 0.5

  const steps = [
    {
      title: 'Upload File',
      description: 'Put your files here.',
      target: () => ref1.current
    }
  ]

  return (
    <div
      ref={containerRef}
      style={{
        transform: `scale(${scaleValue})`,
        transformOrigin: 'top left',
        // 关键：确保容器是相对定位，作为挂载锚点
        position: 'relative',
        width: '200%', // 补偿 scale 造成的空间减小
        height: '200%'
      }}
    >
      <Button ref={ref1}>Upload</Button>

      <Tour
        open={true}
        steps={steps}
        // 关键点 1: 挂载到 transform 的容器内部
        getPopupContainer={() => containerRef.current || document.body}
        // 关键点 2: 因为挂载在内部，Tour 的蒙层和浮层也会被 scale。
        // 如果你不希望它们变小，需要通过样式反向放大
        rootClassName="tour-inside-scale"
      />

      {/* 
         通常还需要配合 CSS（见下文）
         如果是 Antd Tour，它的 Mask 也是 SVG，挂载在内部会自动缩放匹配位置，
         这通常是解决"错位"最直接的方法，副作用是 UI 变小。
      */}
    </div>
  )
}
```

如果不能挂载到内部（例如不想让文字变小），必须挂载到 Body，你需要手动封装一个 `target` 获取器，**手动计算抵消 Scale 后的坐标**。这是一个极端的做法。

**最实用的 Workaround (不用写复杂 JS)**

在 Issues 中被验证有效的方案通常是：**避免在该层级使用 CSS Scale transform**。

如果必须用 Scale（如大屏适配），请使用 **`zoom`** 属性（如果仅针对 Chrome 内核环境开发）。

```css
.your-scale-container {
  /* 替换 transform: scale(0.5) */
  zoom: 0.5;
}
```

这可以完美解决 Popper.js/Floating-ui 的定位问题，因为 `zoom` 改变了布局尺寸，而不只是视觉投影。

### 总结

1.  **首选方案：** 尝试将 `transform: scale` 替换为 CSS `zoom`（如果是后台管理系统，通常只用兼容 Chrome，此法最简单有效，无错位）。
2.  **次选方案：** 如果必须用 `transform: scale`，使用 `getPopupContainer` 将 `Tour` 挂载到产生 `scale` 的那个 DOM 节点内部。
    - _副作用：_ Tour 的 UI（文字、按钮）也会随之缩放。
3.  **终极方案：** 如果上述都不可行，等待 Ant Design 或底层 `rc-tour` 升级支持 `scale` 上下文检测（目前 Antd 正在逐步解决此类问题，但在复杂嵌套 scale 下依然脆弱）。

---

Ant Design 的 `getPopupContainer` 属性主要用于解决浮层（Tooltip, Select 下拉框, Modal 等）渲染位置的问题。

简单来说，它的核心逻辑是：**“找到一个挂载点（DOM 节点），然后利用 React Portal 技术把浮层‘传送’过去。”**

以下是其内部工作流程的逻辑拆解：

### 1. 初始化与默认行为

当一个带有浮层的组件（例如 `<Select>`）渲染时，Antd 会检查是否传入了 `getPopupContainer` 属性。

- **如果未传入（默认）：** 组件通常默认挂载到 `document.body` 上。这意味着浮层的 DOM 结构会直接出现在 `<body>` 标签的直接子级中，独立于当前的 React 组件树结构。

### 2. 触发渲染流程

当组件需要显示浮层时（例如点击了 Select 框），Antd 的底层触发器（通常是 `rc-trigger` 或 `rc-util` 中的 Portal 组件）开始工作。此时，`getPopupContainer` 函数被调用。

### 3. 执行 getPopupContainer 函数

Antd 会调用你传入的这个函数，并将触发浮层显示的那个 DOM 节点（Trigger Node）作为参数传给你。

```javascript
// 伪代码演示内部调用
const triggerNode = this.triggerRef.current // 比如 Select 的输入框 DOM
let mountNode

if (props.getPopupContainer) {
  // 调用用户提供的函数，传入触发节点
  mountNode = props.getPopupContainer(triggerNode)
} else {
  // 默认回退
  mountNode = document.body
}
```

### 4. 确定挂载点（Mount Node）

函数返回的结果就是目标挂载点。Antd 此时拿到了一个真实的 DOM 元素引用。

- **常见用法：** 返回 `triggerNode.parentNode`（父节点），让浮层跟随父容器滚动。

### 5. 创建 React Portal

这是最关键的一步。Antd 使用 `ReactDOM.createPortal` API。

- **逻辑：** 虽然在 React 组件树（Virtual DOM Tree）中，Select 的 Option 列表是 Select 的子组件；但在浏览器真实的 DOM 树中，Antd 会强制把 Option 列表的 DOM 插入到第 4 步找到的 `mountNode` 中。

```javascript
// 伪代码：React Portal 的作用
return ReactDOM.createPortal(
  <PopupContent />, // 浮层的 React 内容
  mountNode // 第 4 步确定的 DOM 容器
)
```

### 6. 定位计算（Positioning）

一旦 DOM 被挂载到指定容器内，Antd 还需要计算浮层的绝对位置（top, left）。

- 它会使用定位库（如 `rc-align` 或 `floating-ui` 等逻辑）计算浮层相对于 **触发节点** 的位置。
- **注意：** 如果你改变了 `getPopupContainer`，通常需要确保该容器拥有 `position: relative` 或 `position: absolute` 属性，否则浮层内部的 `absolute` 定位可能会基于 `body` 偏移，导致错位。

### 总结图示

1.  **用户操作** -> float layer 需要显示。
2.  **Antd 询问** -> "我有 `getPopupContainer` 吗？"
    - **有** -> 运行函数，拿到 DOM 节点 A。
    - **无** -> 使用 `document.body`。
3.  **Portal 传送** -> `ReactDOM.createPortal(浮层, A)`。
4.  **浏览器渲染** -> 浮层的 HTML 源码出现在节点 A 的内部。

### 为什么需要它？

主要用于解决以下问题：

1.  **滚动跟随：** 默认挂载到 body 时，如果你的表格在滚动，浮层可能会“飘”在原处不动。挂载到滚动区域内部可以解决此问题。
2.  **样式隔离/继承：** 挂载到特定容器内，可以让浮层通过 CSS 选择器继承该容器的样式。
3.  **层级上下文（z-index）：** 解决某些父级 `overflow: hidden` 或 `z-index` 导致的遮挡问题。

---

## getPopupContainer 在 Ant Design 组件内部的工作流程

### 1. 整体架构

Ant Design 的弹出层组件都基于 `rc-trigger` 库实现，核心流程：

```
┌─────────────────────────────────────────────────────────────┐
│                    组件层级结构                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Select / Dropdown / Tooltip / Popover / DatePicker        │
│         ↓                                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  rc-trigger (核心触发器库)                           │   │
│  │  - 管理弹出层的显示/隐藏                             │   │
│  │  - 计算弹出层位置                                    │   │
│  │  - 处理 getPopupContainer                           │   │
│  └─────────────────────────────────────────────────────┘   │
│         ↓                                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  rc-align (对齐库)                                   │   │
│  │  - 计算精确的对齐位置                                │   │
│  │  - 处理边界检测和自动调整                            │   │
│  └─────────────────────────────────────────────────────┘   │
│         ↓                                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  React Portal                                        │   │
│  │  - 将弹出层渲染到指定容器                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 2. rc-trigger 内部实现

```typescript
// 简化版 rc-trigger 核心逻辑
class Trigger extends React.Component {
  // 获取弹出层容器
  getContainer = () => {
    const { getPopupContainer, getDocument } = this.props
    const doc = getDocument ? getDocument() : document

    // 1. 如果提供了 getPopupContainer，使用它
    if (getPopupContainer) {
      const triggerNode = this.triggerRef.current
      return getPopupContainer(triggerNode)
    }

    // 2. 否则默认返回 document.body
    return doc.body
  }

  // 渲染弹出层
  renderPopup() {
    const container = this.getContainer()

    // 使用 React Portal 将弹出层渲染到容器中
    return ReactDOM.createPortal(
      <Popup ref={this.popupRef} align={this.getAlign()} onAlign={this.onAlign}>
        {this.props.popup}
      </Popup>,
      container
    )
  }
}
```

### 3. 弹出层位置计算流程

```
┌─────────────────────────────────────────────────────────────┐
│  第一步：获取触发元素的位置                                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  const triggerRect = triggerElement.getBoundingClientRect() │
│                                                             │
│  返回值（相对于视口）:                                       │
│  {                                                          │
│    top: 100,      // 元素顶部到视口顶部的距离                │
│    left: 200,     // 元素左边到视口左边的距离                │
│    width: 120,    // 元素宽度                               │
│    height: 32,    // 元素高度                               │
│    bottom: 132,   // 元素底部到视口顶部的距离                │
│    right: 320     // 元素右边到视口左边的距离                │
│  }                                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  第二步：获取容器元素的位置                                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  const container = getPopupContainer(triggerElement)        │
│  const containerRect = container.getBoundingClientRect()    │
│                                                             │
│  返回值:                                                    │
│  {                                                          │
│    top: 50,       // 容器顶部到视口顶部的距离                │
│    left: 100,     // 容器左边到视口左边的距离                │
│    ...                                                      │
│  }                                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  第三步：计算相对位置                                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  // 弹出层需要定位在容器内部，所以要计算相对位置              │
│  const relativeLeft = triggerRect.left - containerRect.left │
│                     = 200 - 100 = 100                       │
│                                                             │
│  const relativeTop = triggerRect.bottom - containerRect.top │
│                    = 132 - 50 = 82                          │
│                                                             │
│  // 加上容器的滚动偏移                                       │
│  const scrollLeft = container.scrollLeft                    │
│  const scrollTop = container.scrollTop                      │
│                                                             │
│  finalLeft = relativeLeft + scrollLeft                      │
│  finalTop = relativeTop + scrollTop                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  第四步：应用定位样式                                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  popup.style.position = 'absolute'                          │
│  popup.style.left = `${finalLeft}px`                        │
│  popup.style.top = `${finalTop}px`                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4. rc-align 详细计算逻辑

```typescript
// 简化版 rc-align 核心逻辑
function alignElement(source, target, align) {
  // source: 弹出层元素
  // target: 触发元素
  // align: 对齐配置，如 { points: ['tl', 'bl'] }

  // 获取元素位置
  const sourceRect = source.getBoundingClientRect()
  const targetRect = target.getBoundingClientRect()

  // 获取容器（弹出层的定位参照物）
  const container = source.offsetParent || document.body
  const containerRect = container.getBoundingClientRect()

  // 解析对齐点
  // 'tl' = top-left, 'bl' = bottom-left, 'tc' = top-center
  const [sourcePoint, targetPoint] = align.points

  // 计算目标点的绝对位置
  const targetX = getPointX(targetRect, targetPoint)
  const targetY = getPointY(targetRect, targetPoint)

  // 计算源点需要对齐到的位置
  const sourceOffsetX = getPointOffsetX(sourceRect, sourcePoint)
  const sourceOffsetY = getPointOffsetY(sourceRect, sourcePoint)

  // 计算最终位置（相对于容器）
  let left = targetX - containerRect.left - sourceOffsetX
  let top = targetY - containerRect.top - sourceOffsetY

  // 加上滚动偏移
  left += container.scrollLeft
  top += container.scrollTop

  // 应用偏移量
  if (align.offset) {
    left += align.offset[0]
    top += align.offset[1]
  }

  return { left, top }
}

function getPointX(rect, point) {
  switch (
    point[1] // 第二个字符: l=left, c=center, r=right
  ) {
    case 'l':
      return rect.left
    case 'c':
      return rect.left + rect.width / 2
    case 'r':
      return rect.right
  }
}

function getPointY(rect, point) {
  switch (
    point[0] // 第一个字符: t=top, c=center, b=bottom
  ) {
    case 't':
      return rect.top
    case 'c':
      return rect.top + rect.height / 2
    case 'b':
      return rect.bottom
  }
}
```

### 5. 对齐点示意图

```
┌─────────────────────────────────────────────────────────────┐
│  对齐点命名规则: [垂直位置][水平位置]                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│      tl ────── tc ────── tr                                 │
│       │                   │                                 │
│       │                   │                                 │
│      cl        cc        cr     ← 元素                      │
│       │                   │                                 │
│       │                   │                                 │
│      bl ────── bc ────── br                                 │
│                                                             │
│  t = top,    c = center,    b = bottom                      │
│  l = left,   c = center,    r = right                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  常见对齐配置示例                                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 下拉菜单 (Dropdown)                                     │
│     points: ['tl', 'bl']                                    │
│     弹出层的 top-left 对齐到触发器的 bottom-left            │
│                                                             │
│     ┌──────────┐                                            │
│     │  触发器  │ ← bl                                       │
│     └──────────┘                                            │
│     ┌──────────┐                                            │
│  tl→│  弹出层  │                                            │
│     └──────────┘                                            │
│                                                             │
│  2. 向上弹出 (placement="top")                              │
│     points: ['bl', 'tl']                                    │
│     弹出层的 bottom-left 对齐到触发器的 top-left            │
│                                                             │
│     ┌──────────┐                                            │
│     │  弹出层  │                                            │
│  bl→└──────────┘                                            │
│     ┌──────────┐                                            │
│  tl→│  触发器  │                                            │
│     └──────────┘                                            │
│                                                             │
│  3. Tooltip (placement="right")                             │
│     points: ['cl', 'cr']                                    │
│     弹出层的 center-left 对齐到触发器的 center-right        │
│                                                             │
│     ┌──────────┐  ┌─────────────┐                          │
│     │  触发器  │cr cl  Tooltip  │                          │
│     └──────────┘  └─────────────┘                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 6. getPopupContainer 的传递链路

```
┌─────────────────────────────────────────────────────────────┐
│  ConfigProvider 层级传递                                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  <ConfigProvider getPopupContainer={() => container}>       │
│    ↓ (通过 React Context 传递)                              │
│                                                             │
│    <Select>                                                 │
│      ↓ (从 Context 获取)                                    │
│      const { getPopupContainer } = useContext(ConfigContext)│
│      ↓ (传递给内部 Trigger)                                 │
│                                                             │
│      <Trigger getPopupContainer={getPopupContainer}>        │
│        ↓ (调用获取容器)                                     │
│        const container = getPopupContainer(triggerNode)     │
│        ↓ (Portal 渲染)                                      │
│                                                             │
│        ReactDOM.createPortal(<Popup />, container)          │
│      </Trigger>                                             │
│    </Select>                                                │
│  </ConfigProvider>                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 7. 优先级规则

```typescript
// 组件内部的 getPopupContainer 解析逻辑
function resolveGetPopupContainer(props, context) {
  // 优先级从高到低：

  // 1. 组件 props 直接传入的
  if (props.getPopupContainer) {
    return props.getPopupContainer
  }

  // 2. 组件特定的 props（如 Table 的 pagination.getPopupContainer）
  if (props.pagination?.getPopupContainer) {
    return props.pagination.getPopupContainer
  }

  // 3. ConfigProvider 全局配置
  if (context.getPopupContainer) {
    return context.getPopupContainer
  }

  // 4. 默认返回 document.body
  return () => document.body
}
```

### 8. 完整的代码流程示例

```typescript
// Select 组件简化实现
const Select = props => {
  const triggerRef = useRef()
  const { getPopupContainer: contextGetPopupContainer } = useContext(ConfigContext)

  // 合并 getPopupContainer
  const getPopupContainer = props.getPopupContainer || contextGetPopupContainer

  return (
    <div ref={triggerRef} className="select-trigger">
      {/* 触发器内容 */}
      <span>{selectedValue}</span>

      {/* rc-trigger 处理弹出逻辑 */}
      <Trigger
        popup={<DropdownMenu options={options} />}
        getPopupContainer={
          getPopupContainer ? () => getPopupContainer(triggerRef.current) : undefined
        }
        popupAlign={{
          points: ['tl', 'bl'], // 弹出层 top-left 对齐触发器 bottom-left
          offset: [0, 4], // 向下偏移 4px
          overflow: {
            adjustX: true, // 允许自动调整 X 位置
            adjustY: true // 允许自动调整 Y 位置
          }
        }}
      >
        <div className="select-selector" />
      </Trigger>
    </div>
  )
}

// Trigger 内部渲染逻辑
class Trigger extends Component {
  getContainer = () => {
    const { getPopupContainer } = this.props
    if (getPopupContainer) {
      return getPopupContainer()
    }
    return document.body
  }

  render() {
    return (
      <>
        {this.props.children}
        {this.state.visible && (
          <Portal getContainer={this.getContainer}>
            <Align target={() => this.triggerRef.current} align={this.props.popupAlign}>
              {this.props.popup}
            </Align>
          </Portal>
        )}
      </>
    )
  }
}

// Portal 实现
const Portal = ({ children, getContainer }) => {
  const container = useMemo(() => getContainer(), [])
  return ReactDOM.createPortal(children, container)
}
```

### 9. transform: scale 影响分析

```
┌─────────────────────────────────────────────────────────────┐
│  为什么 transform: scale 会导致位置错误？                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  getBoundingClientRect() 返回的是【视觉坐标】               │
│  即经过所有 CSS transform 变换后的最终屏幕位置              │
│                                                             │
│  示例：                                                     │
│  ┌────────────────────────────────────────┐                │
│  │ 画布容器                               │                │
│  │ transform: scale(0.5)                  │                │
│  │ 原始尺寸: 1000x800                     │                │
│  │ 视觉尺寸: 500x400                      │                │
│  │                                        │                │
│  │   ┌──────────┐                         │                │
│  │   │ Select   │                         │                │
│  │   │ 原始位置: (200, 400)               │                │
│  │   │ 视觉位置: (100, 200) ← getBounding │                │
│  │   └──────────┘            ClientRect   │                │
│  └────────────────────────────────────────┘                │
│                                                             │
│  当弹出层渲染到 document.body 时：                          │
│  - 使用视觉坐标 (100, 200) 定位                            │
│  - 但 body 没有 scale 变换                                 │
│  - 弹出层出现在屏幕 (100, 200) 位置                        │
│  - 而视觉上 Select 在缩放后的容器内的 (100, 200)           │
│  - 两者看起来对齐，但如果容器有位移就会错位                 │
│                                                             │
│  更严重的问题：                                              │
│  - 弹出层的尺寸是正常的（未缩放）                           │
│  - 触发器的尺寸是缩放后的                                   │
│  - 导致弹出层看起来比触发器大得多                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

这就是为什么需要使用反向缩放容器 (`useNegativeCameraCss`) 的根本原因——让弹出层的定位和尺寸计算都在同一个坐标系中进行。

---

## getPopupContainer 作用与弹出层定位原理

### 1. 弹出层定位的基本原理

Ant Design/Auxo 等 UI 库的弹出层（Dropdown、Select、Tooltip、Popover 等）定位过程：

```
┌─────────────────────────────────────────────────────────────┐
│  1. 获取触发元素位置                                          │
│     triggerRect = trigger.getBoundingClientRect()           │
│     返回: { top, left, width, height } (相对于视口)          │
│                                                             │
│  2. 获取容器元素位置                                          │
│     containerRect = container.getBoundingClientRect()       │
│                                                             │
│  3. 计算弹出层位置                                           │
│     popupLeft = triggerRect.left - containerRect.left       │
│     popupTop = triggerRect.top - containerRect.top          │
│                                                             │
│  4. 将弹出层渲染到容器中                                      │
│     container.appendChild(popup)                            │
│     popup.style.left = popupLeft                            │
│     popup.style.top = popupTop                              │
└─────────────────────────────────────────────────────────────┘
```

### 2. transform: scale() 导致的问题

```
正常情况 (zoom = 1):
┌──────────────────────────────────────────┐
│ document.body                            │
│  ┌────────────────────────────────────┐  │
│  │ 画布容器 (transform: scale(1))     │  │
│  │  ┌──────────┐                      │  │
│  │  │ Select   │ ← getBoundingClientRect() │
│  │  │ 触发器   │   返回: left=100, top=200  │
│  │  └──────────┘                      │  │
│  │                                    │  │
│  └────────────────────────────────────┘  │
│  ┌──────────┐                            │
│  │ 下拉菜单 │ ← 渲染到 body              │
│  │ left=100 │   定位: left=100, top=220  │
│  │ top=220  │   ✅ 位置正确！            │
│  └──────────┘                            │
└──────────────────────────────────────────┘

缩放情况 (zoom = 0.5):
┌──────────────────────────────────────────┐
│ document.body                            │
│  ┌─────────────────────┐                 │
│  │ 画布 scale(0.5)     │                 │
│  │  ┌─────┐            │                 │
│  │  │Select│ ← getBoundingClientRect()   │
│  │  └─────┘   返回: left=50, top=100     │
│  │            (实际视觉位置，被缩放了)    │
│  └─────────────────────┘                 │
│                                          │
│        ┌──────────┐                      │
│        │ 下拉菜单 │ ← 渲染到 body        │
│        │ left=50  │   定位基于缩放后坐标 │
│        │ top=120  │   ❌ 位置偏左上！    │
│        └──────────┘                      │
└──────────────────────────────────────────┘
```

**关键问题**：`getBoundingClientRect()` 返回的是**经过 transform 变换后的视觉坐标**，但弹出层渲染在未缩放的 `document.body` 中。

### 3. getPopupContainer 的作用

```typescript
// Ant Design Select 组件示例
<Select
  getPopupContainer={(triggerNode) => triggerNode.parentElement}
/>

// ConfigProvider 全局配置
<ConfigProvider
  getPopupContainer={() => document.querySelector('#my-container')}
>
```

**作用**：指定弹出层的挂载容器，影响：

1. **DOM 位置**：弹出层渲染到哪个元素内
2. **坐标计算**：弹出层的 left/top 相对于这个容器计算
3. **层叠上下文**：弹出层的 z-index 在这个容器的层叠上下文中生效

### 4. 解决方案的原理

```
解决方案：反向缩放容器
┌──────────────────────────────────────────┐
│ document.body                            │
│  ┌──────────────────────────────────┐    │
│  │ .tl-container                    │    │
│  │  ┌────────────────────────────┐  │    │
│  │  │ .tl-layer scale(0.5)       │  │    │
│  │  │  ┌──────────┐              │  │    │
│  │  │  │ Select   │              │  │    │
│  │  │  │ 触发器   │              │  │    │
│  │  │  └──────────┘              │  │    │
│  │  └────────────────────────────┘  │    │
│  │                                  │    │
│  │  ┌────────────────────────────┐  │    │
│  │  │ .tl-popup scale(2)        │  │    │
│  │  │ (反向缩放: 1/0.5 = 2)     │  │    │
│  │  │  ┌──────────┐              │  │    │
│  │  │  │ 下拉菜单 │ ← 渲染到这里 │  │    │
│  │  │  │ 正确对齐 │              │  │    │
│  │  │  └──────────┘              │  │    │
│  │  └────────────────────────────┘  │    │
│  └──────────────────────────────────┘    │
└──────────────────────────────────────────┘
```

**原理**：

1. 弹出层容器与画布层在同一父容器内
2. 弹出层容器应用 `transform: scale(1/zoom)`（反向缩放）
3. 弹出层渲染到反向缩放容器中
4. 最终视觉效果：弹出层与触发器正确对齐

### 5. 坐标计算详解

```javascript
// 假设 zoom = 0.5
// 触发器在画布内的实际位置: (200, 400)
// 画布缩放后的视觉位置: (100, 200)

// 情况 A: 弹出层渲染到 document.body
triggerRect = { left: 100, top: 200 } // 缩放后的视觉坐标
containerRect = { left: 0, top: 0 } // body 的位置
popupPosition = { left: 100, top: 220 } // ❌ 错误！相对于未缩放的空间

// 情况 B: 弹出层渲染到反向缩放容器 (scale(2))
// 反向缩放容器的 getBoundingClientRect 也会受到外层 scale(0.5) 影响
// 但内部的弹出层会被 scale(2) 放大，最终抵消

// 更精确的方案：使用 CSS transform-origin 和精确计算
```

### 6. useNegativeCameraCss 的实现

```typescript
// packages/disco-core/src/hooks/useCameraCss.tsx
export function useNegativeCameraCss(ref, pageState) {
  React.useLayoutEffect(() => {
    const { zoom, point } = pageState.camera
    const elm = ref.current
    if (!elm) return

    // 反向缩放：1/zoom
    // 反向平移：抵消画布的平移
    elm.style.setProperty(
      'transform',
      `scale(${1 / zoom}) translateX(${-point[0]}px) translateY(${-point[1]}px)`
    )
  }, [pageState.camera.zoom, pageState.camera.point])
}
```

### 7. 为什么 Table/Pagination 的情况特殊

```
组件嵌套结构：
ConfigProvider (getPopupContainer = () => body)
  └── 画布
       └── Table
            └── Pagination
                 └── Select
                      └── 内部自己设置了 getPopupContainer
                          └── 弹出层渲染到 Select 的父元素
```

**问题**：Table/Pagination 组件内部可能覆盖了全局的 `getPopupContainer` 配置，导致弹出层没有渲染到我们指定的反向缩放容器中。

**解决方案**：需要在组件渲染层面提供 `getPopupContainer` 方法，让画布内的所有组件都能获取到正确的弹出层容器。

---
