ProseMirror 的原生 API 非常强大但极其琐碎。**TipTap** 的成功在于它将 ProseMirror 的 `Schema`、`Plugins`、`Keymaps` 和 `InputRules` 封装进了一个高度模块化的 **Extension (扩展)** 系统。

模仿 TipTap 的架构，我们可以将编辑器拆分为：**Core (核心管理器)** 和 **Extensions (功能插件)**。

### 1. 推荐的目录结构

这种结构实现了功能的高度内聚。如果你想增加“表格”功能，只需在 `extensions` 下增加一个文件夹，而无需修改核心逻辑。

```text
src/editor/
├── core/
│   ├── extension.ts          # 扩展基类定义
│   ├── extension-manager.ts  # 负责合并 Schema 和 Plugins
│   └── editor.ts             # 暴露给 UI 框架的入口类
├── extensions/
│   ├── base/                 # 基础插件 (History, Keymap)
│   ├── nodes/                # 节点扩展 (Paragraph, Heading, Image)
│   │   ├── paragraph.ts
│   │   └── heading.ts
│   └── marks/                # 标记扩展 (Bold, Italic, Link)
│       ├── bold.ts
│       └── link.ts
└── types/                    # 类型定义
```

### 2. 核心实现：Extension 基类

这是模仿 TipTap 的关键。每个扩展决定了它如何影响 Schema 和编辑器行为。

```typescript
// src/editor/core/extension.ts
import { NodeSpec, MarkSpec } from 'prosemirror-model'
import { Plugin } from 'prosemirror-state'

export abstract class Extension {
  abstract name: string

  // 贡献节点定义
  addOptions() {
    return {}
  }

  // 贡献 Schema
  addNodeSpec?(): NodeSpec
  addMarkSpec?(): MarkSpec

  // 贡献插件 (Keymap, InputRules 等)
  addProseMirrorPlugins?(): Plugin[]

  // 贡献命令
  addCommands?(): Record<string, any>
}
```

### 3. 核心实现：Editor 管理器

管理器负责遍历所有扩展，动态生成 `Schema` 并初始化 `EditorView`。

```typescript
// src/editor/core/editor.ts
import { Schema } from 'prosemirror-model'
import { EditorState } from 'prosemirror-state'
import { EditorView } from 'prosemirror-view'
import { Extension } from './extension'

export class MyEditor {
  public schema: Schema
  public view: EditorView

  constructor(options: { extensions: Extension[]; content: string; element: HTMLElement }) {
    // 1. 从扩展中提取 Nodes 和 Marks
    const nodes: any = {}
    const marks: any = {}

    options.extensions.forEach(ext => {
      if (ext.addNodeSpec) nodes[ext.name] = ext.addNodeSpec()
      if (ext.addMarkSpec) marks[ext.name] = ext.addMarkSpec()
    })

    // 2. 创建 Schema
    this.schema = new Schema({ nodes, marks })

    // 3. 合并所有插件
    const plugins = options.extensions.flatMap(ext =>
      ext.addProseMirrorPlugins ? ext.addProseMirrorPlugins() : []
    )

    // 4. 初始化 View
    this.view = new EditorView(options.element, {
      state: EditorState.create({
        schema: this.schema,
        plugins
      })
    })
  }
}
```

### 4. 扩展实例：Bold (Mark)

通过这种方式，Bold 的逻辑（HTML 标签、解析规则、快捷键）被完美封装。

```typescript
// src/editor/extensions/marks/bold.ts
import { Extension } from '../../core/extension'
import { toggleMark } from 'prosemirror-commands'
import { keymap } from 'prosemirror-keymap'

export class Bold extends Extension {
  name = 'bold'

  addMarkSpec() {
    return {
      parseDOM: [
        { tag: 'strong' },
        { tag: 'b', getAttrs: (node: any) => node.style.fontWeight != 'normal' && null },
        {
          style: 'font-weight',
          getAttrs: (value: string) => /^(bold(er)?|[5-9]\d{2,})$/.test(value) && null
        }
      ],
      toDOM: () => ['strong', 0]
    }
  }

  addProseMirrorPlugins() {
    return [
      keymap({
        'Mod-b': toggleMark(this.editor.schema.marks.bold)
      })
    ]
  }
}
```

### 5. 为什么这是最佳实践？

1.  **Schema 动态化**：在原生 ProseMirror 中，Schema 是静态定义的。这种架构允许你根据用户权限或业务场景，动态传入不同的 `extensions` 数组，从而生成不同的 Schema。
2.  **命令封装**：你可以像 TipTap 一样实现链式调用 `editor.commands.toggleBold().run()`，只需在管理器中对 `dispatch` 进行包装。
3.  **NodeView 友好**：可以在 Extension 中定义 `addNodeView` 方法，将复杂的 React/Vue 组件与特定的 Node 绑定。
4.  **易于测试**：每个 Extension 都是独立的类，可以针对单个功能编写单元测试。

这种架构将 ProseMirror 从一个“库”提升为了一个“平台”，非常适合需要长期维护的复杂富文本项目。
