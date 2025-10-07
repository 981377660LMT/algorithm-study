好的，我们来详细讲解 `tldraw` 中关于 **Shapes (图形)** 的概念。这篇文档是 `tldraw` 定制和扩展的核心，它解释了图形是什么、如何工作以及如何创建你自己的自定义图形。

---

### **什么是图形 (Shape)？**

在 `tldraw` 中，一个**图形 (Shape)** 是指任何可以存在于画布页面上的可见元素，例如箭头、图片、文本框或几何图形。

### **图形的类型**

`tldraw` 将图形分为三种类型：

1.  **核心图形 (Core Shapes)**

    - 这些是编辑器内置且始终存在的图形。它们是编辑器基础功能的一部分。
    - 目前，唯一的**核心图形**是**组合 (Group)** 图形。

2.  **默认图形 (Default Shapes)**

    - 这些是 `<Tldraw />` 组件默认包含的所有图形，例如箭头 (`TLArrowShape`)、手绘 (`TLDrawShape`)、几何图形 (`TLGeoShape`) 等。
    - 你可以从 `tldraw` 库中通过 `defaultShapeUtils` 导入它们。

3.  **自定义图形 (Custom Shapes)**
    - 这些是由你（开发者）自己创建的图形，用于扩展 `tldraw` 的功能。本文档的重点就是如何创建自定义图形。

---

### **图形对象 (The Shape Object) - “数据”**

在 `tldraw` 的数据模型中，每个图形都只是一个存储在 `store` 中的 **JSON 对象**（官方称为 `record`）。

例如，这是一个矩形 (`geo`) 图形的数据结构：

```json
{
  "parentId": "page:somePage",
  "id": "shape:someId",
  "typeName": "shape",
  "type": "geo",
  "x": 106,
  "y": 294,
  "rotation": 0,
  "index": "a28",
  "opacity": 1,
  "isLocked": false,
  "props": {
    "w": 200,
    "h": 200,
    "geo": "rectangle",
    "color": "black",
    "fill": "none",
    "dash": "draw",
    "size": "m",
    "text": "diagram"
    // ... 其他 geo 图形特有的属性
  },
  "meta": {}
}
```

这个 JSON 对象包含三个关键部分：

1.  **基础属性 (Base Properties)**

    - 所有图形都共享一些基础信息，如 `id`、`type` (图形类型名)、位置 (`x`, `y`)、旋转 (`rotation`)、层级 (`index`)、父级 (`parentId`) 等。

2.  **`props` (特有属性)**

    - 这部分包含了**特定类型**图形独有的信息。不同类型的图形 `props` 也完全不同。
    - 例如，文本图形的 `props` 包含文本内容，而箭头图形的 `props` 则包含端点样式等。

3.  **`meta` (元数据)**
    - 这是一个“逃生舱口”，用于存储你的应用需要但 `tldraw` 本身不关心的额外数据。
    - 例如，你可以用它来存储创建者的用户 ID、创建日期或其他业务相关的信息。

---

### **`ShapeUtil` 类 - “逻辑”**

如果说图形对象是“数据”，那么 `ShapeUtil` 类就是定义该图形所有**“行为和逻辑”**的地方。

`tldraw` 的编辑器本身并不知道如何渲染一个矩形或处理它的交互。当编辑器需要对一个图形进行操作时（例如渲染、计算边界、响应缩放），它会找到该图形 `type` 对应的 `ShapeUtil` 类，并调用其方法来获取答案。

例如，当编辑器需要渲染一个 `type` 为 `'card'` 的图形时，它会找到 `CardShapeUtil` 并调用其 `component` 方法。

### **如何创建自定义图形**

下面我们通过创建一个简单的 "card" 图形（一个带文字的矩形）来了解完整流程。

#### **第 1 步：定义图形的类型 (Shape Type)**

首先，使用 TypeScript 定义你的图形对象的结构。`TLBaseShape` 是一个辅助类型。

```typescript
import { TLBaseShape } from 'tldraw'

// 定义一个名为 'card' 的图形类型
// 它的 props 包含宽度 (w) 和高度 (h)
type CardShape = TLBaseShape<'card', { w: number; h: number }>
```

`TLBaseShape` 会自动为你添加 `id`, `x`, `y` 等基础属性。

#### **第 2 步：创建 `ShapeUtil` 类**

接下来，创建 `ShapeUtil` 类来定义图形的行为。

```typescript
import { HTMLContainer, ShapeUtil, Rectangle2d } from 'tldraw'

class CardShapeUtil extends ShapeUtil<CardShape> {
  // 1. 静态 type 属性，必须与你的图形类型名匹配
  static override type = 'card' as const

  // 2. 定义图形的默认 props
  getDefaultProps(): CardShape['props'] {
    return {
      w: 100,
      h: 100
    }
  }

  // 3. 定义图形的几何形状，用于碰撞检测和绑定
  getGeometry(shape: CardShape) {
    return new Rectangle2d({
      width: shape.props.w,
      height: shape.props.h,
      isFilled: true
    })
  }

  // 4. 定义如何渲染图形的 React 组件
  component(shape: CardShape) {
    // HTMLContainer 是一个辅助组件，用于包裹 DOM 元素
    return (
      <HTMLContainer>
        <div
          style={{
            width: shape.props.w,
            height: shape.props.h,
            border: '1px solid black',
            textAlign: 'center'
          }}
        >
          Hello
        </div>
      </HTMLContainer>
    )
  }

  // 5. 定义图形被选中时显示的指示器（蓝色边框）
  indicator(shape: CardShape) {
    return <rect width={shape.props.w} height={shape.props.h} />
  }
}
```

这是一个最小化的 `ShapeUtil` 实现，它定义了图形的类型、默认属性、几何边界、渲染组件和选中指示器。

#### **第 3 步：将 `ShapeUtil` 注入 `Tldraw` 组件**

将你创建的 `ShapeUtil` 数组通过 `shapeUtils` prop 传递给 `<Tldraw>` 组件。

```tsx
const MyCustomShapes = [CardShapeUtil]

export default function App() {
  return (
    <div style={{ position: 'fixed', inset: 0 }}>
      <Tldraw
        shapeUtils={MyCustomShapes}
        onMount={editor => {
          // 在画布加载后，创建一个 'card' 类型的图形
          editor.createShapes([{ type: 'card' }])
        }}
      />
    </div>
  )
}
```

刷新页面后，你就能在画布上看到你创建的自定义图形了。

---

### **高级功能**

- **文本标签 (Labels)**: `tldraw` 内置了对纯文本和富文本（使用 [TipTap](https://tiptap.dev/)）的支持。你可以使用 `PlainTextLabel` 或 `RichTextLabel` 组件轻松地为你的自定义图形添加可编辑的文本。
- **交互 (Interaction)**: 你可以设置 `ShapeUtil` 的标志（flags）来启用图形内部的指针事件，允许用户与图形内容直接交互。
- **裁剪 (Clipping)**: 自定义图形可以通过实现 `getClipPath` 和 `shouldClipChild` 方法来裁剪其子图形。
- **迁移 (Migrations)**: 如果你的图形 `props` 结构发生了变化，你可以为 `ShapeUtil` 添加 `migrations` 来确保旧数据可以平滑地升级到新结构。
- **继承基础行为**: 你可以继承 `BaseBoxShapeUtil` 等“起始”工具类，来快速获得标准的矩形缩放、旋转等行为。
