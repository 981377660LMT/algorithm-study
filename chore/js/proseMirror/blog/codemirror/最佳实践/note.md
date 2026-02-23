```
/* eslint-disable @typescript-eslint/ban-ts-comment */
/* eslint-disable arrow-body-style */
import {
  Decoration,
  DecorationSet,
  EditorView,
  MatchDecorator,
  ViewPlugin,
  ViewUpdate
} from '@codemirror/view';
import { Extension } from '@codemirror/state';
import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget';
import type { CapsuleGlobalState } from '../types';
import { CAPSULE_REGEX, extractAllCapsules } from '../utils';
import type { ICapsuleInfo } from '@/components/ChatProvider';

// 生成稳定的 key 用于比较胶囊变化
const getCapsulesKey = (capsules: ICapsuleInfo[]): string =>
  JSON.stringify(capsules.map((c) => c.id));

/**
 * 创建胶囊插件
 * 遵循 CodeMirror 最佳实践：
 * 1. ViewPlugin 仅负责装饰器的创建和更新
 * 2. updateListener 处理副作用（调用外部回调）
 */
export const createCapsulePlugin = (globalState: CapsuleGlobalState): Extension => {
  const matcher = new MatchDecorator({
    regexp: CAPSULE_REGEX,
    decoration: (match) => {
      const name = match[1];
      if (!name || name.length === 0) return null;
      return Decoration.replace({
        widget: new MentionCapsuleWidget(name, globalState),
        inclusive: false
      });
    }
  });

  // 装饰器插件：仅负责渲染，不处理副作用
  const decorationPlugin = ViewPlugin.fromClass(
    class CapsuleDecorationPlugin {
      decorations: DecorationSet;

      constructor(view: EditorView) {
        this.decorations = matcher.createDeco(view);
      }

      update(update: ViewUpdate) {
        this.decorations = matcher.updateDeco(update, this.decorations);
      }
    },
    {
      decorations: (plugin) => plugin.decorations,
      provide(plugin) {
        return EditorView.atomicRanges.of((view) => view.plugin(plugin)?.decorations ?? Decoration.none);
      }
    }
  );

  let lastCapsulesKey = '';

  // 更新监听器：在事务完成后检测变化并触发回调
  // 这是处理副作用的正确位置，而不是在 ViewPlugin.update 中
  const updateListener = EditorView.updateListener.of((update) => {
    if (!update.docChanged) return;

    const text = update.state.doc.toString();
    const allCommands = globalState.commandsGetter?.() || [];
    const capsules = extractAllCapsules(text, allCommands);
    const capsulesKey = getCapsulesKey(capsules);

    if (capsulesKey !== lastCapsulesKey) {
      lastCapsulesKey = capsulesKey;
      // 使用 queueMicrotask 确保回调在当前事务完全完成后执行
      queueMicrotask(() => {
        globalState.capsulesChangeHandler?.(capsules);
      });
    }
  });

  return [decorationPlugin, updateListener];
};

```

已进行更改。

我已经完成了代码的 Review 和重构。针对您提出的两个问题，以下是详细的分析和修改说明：

### 1. 代码重构：遵循 CodeMirror 最佳实践

我重构了 capsule-plugin.ts，主要解决了**副作用处理不当**的问题。

**修改前的问题：**

- 在 `ViewPlugin` 的 `update` 方法中直接执行业务逻辑（`checkCapsuleChange`）并触发 React 回调。
- **风险**：`ViewPlugin.update` 应仅用于更新视图（如 Decorations）。在此处触发外部状态更新（React `setState`）可能导致 "Render Loop" 或时序冲突，这是 CodeMirror 开发中的常见反模式。

**修改后的方案（最佳实践）：**

- **分离职责**：
  - `ViewPlugin`：现在只负责 UI 渲染（维护 `DecorationSet`），保持纯净。
  - `EditorView.updateListener`：专门用于监听编辑器状态变化。我将 `checkCapsuleChange` 的逻辑移到了这里。
- **安全的副作用调度**：
  - 使用了 `queueMicrotask` 包裹回调函数。这确保了 React 的状态更新在当前 CodeMirror 事务完全结束后才执行，避免了 DOM 更新冲突。
- **类型安全**：修复了相关的 TypeScript 类型定义。

### 2. `checkCapsuleChange` 的方式是否合理？

**结论：合理，且在当前场景下是最佳选择。**

- **当前实现方式**：

  - 每次文档变化（`docChanged`）时，获取全文内容。
  - 使用正则（`extractAllCapsules`）全量扫描提取所有胶囊。
  - 生成指纹（`capsulesKey`）比对是否变化，变化时触发回调。

- **为什么是合理的？**
  - **单一数据源（Single Source of Truth）**：胶囊的存在完全由文本内容（`@name` + 零宽空格）决定。这种方式永远不会出现“内部状态”与“文档内容”不同步的 Bug。
  - **性能无忧**：对于聊天输入框场景（通常文本在几千字以内），现代浏览器的正则引擎全量扫描耗时极短（微秒级），不会造成任何感知的卡顿。
  - **维护简单**：相比于使用 `StateField` 编写复杂的增量更新逻辑（处理插入、删除、撤销重做等各种边界情况），全量扫描的代码极其简洁且健壮。

**修改后的代码摘要：**

```typescript
// ...existing code...
// 装饰器插件：仅负责渲染，不处理副作用
const decorationPlugin = ViewPlugin.fromClass(
  class CapsuleDecorationPlugin {
    // ... 只处理 decorations
  }
  // ...
)

// 更新监听器：在事务完成后检测变化并触发回调
// 这是处理副作用的正确位置
const updateListener = EditorView.updateListener.of(update => {
  if (!update.docChanged) return

  const text = update.state.doc.toString()
  // ... 提取胶囊逻辑 ...

  if (capsulesKey !== lastCapsulesKey) {
    lastCapsulesKey = capsulesKey
    // 使用 queueMicrotask 确保回调在当前事务完全完成后执行
    queueMicrotask(() => {
      globalState.capsulesChangeHandler?.(capsules)
    })
  }
})

return [decorationPlugin, updateListener]
// ...existing code...
```

这次重构使代码更符合 CodeMirror 的设计哲学，提高了可维护性和稳定性。

---

### 评价 `extractAllCapsules` 的方式

当前的实现方式是**基于文本模式的动态解析（Text-based Pattern Matching）**。

#### ✅ 优点

1.  **单一数据源（Single Source of Truth）**：胶囊的存在完全由文本内容（`@name` + 零宽空格）决定。这意味着你永远不需要担心“内部状态”与“编辑器文本”不同步的问题。
2.  **实现简单（KISS）**：逻辑非常直观，无状态，易于调试。
3.  **适合短文本**：在聊天输入框场景下（通常 < 5000 字），正则全量扫描的性能开销几乎可以忽略不计。

#### ❌ 缺点/局限性

1.  **无法处理“同名不同 ID”**：这是最大的硬伤。如果 `commands` 列表中有两个叫 "张三" 的人（ID 不同），当前的逻辑 `findCommandItemByName` 只能找到第一个，导致数据错误。
2.  **强依赖外部数据**：解析结果依赖 `commands` 列表。如果后端更新了列表，移除了某个项，编辑器里的胶囊就会“失效”或匹配错误。
3.  **重命名困难**：无法支持“修改胶囊显示文本但保持 ID 不变”的需求。
4.  **性能隐患**：如果未来扩展到长文档编辑，每次按键都进行全量正则扫描会成为性能瓶颈。

---

### 🚀 更好的方式：基于 StateField 的元数据持久化

如果你的应用需要更健壮的胶囊系统（特别是需要区分同名项），最佳实践是使用 **CodeMirror 的 `StateField` 维护带有元数据的 `DecorationSet`**。

**核心思想**：
不再在读取时通过名字反查 ID，而是在**插入胶囊时**就将 ID 绑定在文档的特定范围（Range）上。

#### 实现思路

1.  **定义胶囊数据结构**：
    在 `Decoration` 的 `spec` 中存储完整的胶囊数据（ID, name, data）。

2.  **使用 StateField 维护状态**：
    创建一个 `StateField`，它维护一个 `RangeSet<Decoration>`。
    - **插入时**：通过 `Transaction` 发送一个 Effect，携带胶囊数据，StateField 接收后在指定位置添加 Range。
    - **更新时**：利用 `RangeSet.map(tr.changes)` 自动处理文本变更带来的位置偏移。CodeMirror 会自动帮你维护胶囊的位置。
    - **读取时**：直接遍历 `state.field(capsuleField)` 获取所有胶囊数据，**无需正则扫描，无需反查列表**。

#### 代码示例

我为你写了一个基于 `StateField` 的改进方案示例。你可以对比一下：

```typescript
import { StateField, StateEffect } from '@codemirror/state'
import { Decoration, DecorationSet, EditorView, WidgetType } from '@codemirror/view'
import { RangeSetBuilder } from '@codemirror/state'

// 1. 定义胶囊数据接口
interface CapsuleData {
  id: string
  name: string
  // 可以包含更多数据，如 type, avatar 等
}

// 2. 定义 Effect 用于添加胶囊
export const addCapsuleEffect = StateEffect.define<{
  from: number
  to: number
  data: CapsuleData
}>()

// 3. 定义 Widget
class CapsuleWidget extends WidgetType {
  constructor(readonly data: CapsuleData) {
    super()
  }
  toDOM() {
    const span = document.createElement('span')
    span.textContent = `@${this.data.name}`
    span.className = 'mention-capsule'
    // ... 渲染逻辑
    return span
  }
}

// 4. 核心：StateField 管理胶囊状态
export const capsuleField = StateField.define<DecorationSet>({
  create() {
    return Decoration.none
  },
  update(capsules, tr) {
    // A. 自动映射位置变化 (核心优势：CodeMirror 帮你处理位置偏移)
    capsules = capsules.map(tr.changes)

    // B. 处理添加胶囊的 Effect
    for (const e of tr.effects) {
      if (e.is(addCapsuleEffect)) {
        const deco = Decoration.widget({
          widget: new CapsuleWidget(e.value.data),
          side: 1,
          // 将数据直接绑定在 Decoration 的 spec 上
          spec: { capsuleData: e.value.data }
        })
        capsules = capsules.update({
          add: [deco.range(e.value.from, e.value.to)] // 这里的 from/to 通常是零宽范围
        })
      }
    }
    return capsules
  },
  provide: f => EditorView.decorations.from(f)
})

// 5. 更好的 extractAllCapsules 实现
export const getCapsulesFromState = (state: EditorState) => {
  const capsules: CapsuleData[] = []
  const rangeSet = state.field(capsuleField)

  // 直接遍历 RangeSet，无需正则，无需反查
  rangeSet.between(0, state.doc.length, (from, to, value) => {
    if (value.spec.capsuleData) {
      capsules.push(value.spec.capsuleData)
    }
  })

  return capsules
}
```

### 总结建议

- **保持现状**：如果你的场景仅仅是简单的 `@Bot` 引用，且名字唯一，当前的正则方式是**性价比最高**的。
- **采用 StateField**：如果你发现有**同名用户区分**的需求，或者需要支持**重命名**，那么请切换到上述的 `StateField` 方案。这虽然增加了代码复杂度，但提供了真正的“数据驱动”能力。

---

您提到的“解析视图”应该是指当前代码中使用的 **基于文本正则解析（Text Parsing / Scanning）** 来生成胶囊视图的方式。

我们将这种方式与 **基于 StateField（状态字段）** 的方式进行深度对比：

### 1. 核心区别

| 特性         | 文本解析模式 (当前方案)         | StateField 模式 (进阶方案)           |
| :----------- | :------------------------------ | :----------------------------------- |
| **数据来源** | **纯文本** (`@name` + 零宽空格) | **内存状态** (StateField 中的元数据) |
| **身份识别** | 靠名字反查 (Name -> ID)         | 直接存储 ID (ID 绑定在 Range 上)     |
| **持久化**   | 文本即数据，天然持久化          | 需要额外的序列化/反序列化逻辑        |
| **复杂度**   | 低 (KISS 原则)                  | 高 (需要管理 Transaction 和 Effect)  |

---

### 2. 文本解析模式 (当前方案)

这是目前代码中的实现方式：`extractAllCapsules` + `MatchDecorator`。

#### ✅ 优点

1.  **单一事实来源 (Single Source of Truth)**：
    - 胶囊的存在完全由文本决定。你永远不需要担心“状态”和“文本”不同步。
    - **所见即所得**：如果用户删除了 `@` 符号，胶囊自动消失，不需要写任何额外的状态清理逻辑。
2.  **复制粘贴极其友好**：
    - 用户从记事本复制一段 `@张三\u200B` 进来，编辑器会自动将其渲染为胶囊，无需任何特殊处理。
3.  **实现极其简单**：
    - 核心逻辑就是一个正则表达式。代码量少，Bug 少，易于维护。
4.  **后端存储简单**：
    - 后端只需要存纯文本字符串。前端加载时自动解析渲染。

#### ❌ 缺点

1.  **无法处理同名异义**：
    - 这是**致命伤**。如果有两个用户都叫“张三”（ID: 001, ID: 002），文本解析模式只能通过名字“张三”去列表里查，通常只能查到第一个。导致 ID: 002 的张三永远无法被正确选中。
2.  **强依赖外部数据**：
    - 渲染必须依赖 `commands` (用户列表)。如果网络慢，列表没加载回来，编辑器里的胶囊就无法渲染或显示为纯文本。
3.  **重命名困难**：
    - 如果想把胶囊显示的文本从“张三”改成“法外狂徒”，你必须修改编辑器里的文本内容，这可能会破坏正则匹配。
4.  **性能隐患 (长文档)**：
    - 每次按键都会触发全量正则扫描。在几万字的长文档中，这会造成输入卡顿（虽然在聊天框场景通常不是问题）。

---

### 3. StateField 模式 (进阶方案)

这是 CodeMirror 官方推荐的复杂富文本场景实现方式。

#### ✅ 优点

1.  **完美支持同名异义**：
    - 胶囊的数据结构是 `{ id: '002', name: '张三' }` 绑定在文档的 `[5, 8]` 坐标上。
    - 即使文本是“张三”，我们也能明确知道它是 ID 002。
2.  **数据独立性**：
    - 胶囊携带了完整的元数据（头像、类型等）。即使外部的用户列表接口挂了，编辑器依然能完美展示胶囊信息。
3.  **高性能 (增量更新)**：
    - CodeMirror 引擎会自动处理位置映射（Mapping）。你在胶囊前面打字，胶囊的坐标会自动后移，不需要重新扫描全文。
4.  **灵活的交互**：
    - 可以轻松实现“点击胶囊弹窗修改属性”、“修改胶囊显示文本但保持 ID 不变”等高级功能。

#### ❌ 缺点

1.  **序列化/反序列化复杂**：
    - **复制粘贴是噩梦**：如果用户复制了胶囊粘贴到另一个地方，默认只能粘贴纯文本。你必须实现自定义的 `Clipboard` 处理逻辑，将元数据序列化为 JSON 或 HTML 属性，粘贴时再解析回来。
    - **后端存储麻烦**：保存到后端时，不能只存纯文本，通常需要存 JSON 结构（如 Delta 格式）或者带标记的 HTML，否则下次加载时元数据就丢了。
2.  **状态同步心智负担**：
    - 需要时刻注意 StateField 中的 Range 是否与文档内容匹配。虽然 CodeMirror 做了很多工作，但开发复杂度依然是文本解析模式的 3-5 倍。

---

### 总结与建议

**什么时候坚持用“文本解析模式”？**

- **聊天输入框**：文本短，生命周期短。
- **名字唯一**：业务上能保证 `@` 的对象名字不重复，或者重复了也不影响核心逻辑。
- **轻量级**：不需要复杂的富文本交互。

**什么时候切换到“StateField 模式”？**

- **文档编辑器**：类似 Notion、飞书文档，需要长期保存和复杂编辑。
- **同名区分**：必须严格区分同名用户。
- **富交互**：胶囊不仅仅是个标签，还是个可交互的对象（比如点击弹出详情卡片，且详情数据不依赖外部列表）。

**结论**：
在您的 **Chat Bot SDK 输入框** 场景下，目前的 **文本解析模式** 是性价比最高的选择。除非您遇到了“同名用户导致 @ 错误”的用户投诉，否则不建议过早优化引入 StateField 的复杂度。

---

### 1. 为什么需要 `event.preventDefault()`？

这两者控制的对象不同：

- **`return true`**：是告诉 **CodeMirror** 编辑器框架 —— “这个事件我已经处理了，请**不要**执行 CodeMirror 内置的粘贴逻辑”。
- **`event.preventDefault()`**：是告诉 **浏览器** —— “请**不要**执行原生的粘贴行为（即不要直接把文本插入到 DOM 节点中）”。

**为什么要同时用？**
如果你只返回 `true` 但不调用 `preventDefault()`：

1.  CodeMirror 确实停手了。
2.  但浏览器不知道，浏览器会继续把剪贴板的内容直接塞进编辑器对应的 DOM 元素里。
3.  这会导致 **DOM 内容** 和 **CodeMirror 的内部状态 (State)** 不一致，引发编辑器崩溃或显示错乱。

### 2. 之前的代码问题

之前的代码问题在于**过早地**调用了 `event.preventDefault()`。

```typescript
// 旧逻辑
event.preventDefault() // 1. 先把浏览器默认行为禁了
if (openRef.current) {
  if (!text.includes(CAPSULE_BOUNDARY)) return false // 2. 然后告诉 CodeMirror "我没处理，你来"
}
```

虽然理论上 CodeMirror 接手后也能处理，但在某些情况下，这种“先禁止原生行为，再让框架处理”的混合状态会导致异常（例如 CodeMirror 可能会检查 `defaultPrevented` 状态，或者浏览器的剪贴板事件流被干扰）。

**最佳实践是：** 只有当你确定**真的要自己处理**这个事件时，才调用 `event.preventDefault()` 并返回 `true`。如果你不想处理，就什么都别做，直接返回 `false`。

我将为您修正这段代码，将 `preventDefault` 移到判断逻辑之后。

```
    const handlePaste = useCallback((event: ClipboardEvent, view: EditorView) => {
      const text = event.clipboardData?.getData('text/plain');
      if (!text) return false;

      if (openRef.current) {
        if (!text.includes(CAPSULE_BOUNDARY)) return false;

        event.preventDefault();
        const cleanText = text.replace(new RegExp(CAPSULE_BOUNDARY, 'g'), '');
        view.dispatch(view.state.replaceSelection(cleanText));
        return true;
      }

      event.preventDefault();
      const regex = new RegExp(AUTO_COMPLETE_CAPSULE_REGEX.source, 'g');
      const processedText = text.replace(regex, `$&${CAPSULE_BOUNDARY}`);
      view.dispatch(view.state.replaceSelection(processedText));
      return true;
    }, []);
```

已进行更改。

```
    const handlePaste = useCallback((event: ClipboardEvent, view: EditorView) => {
      const text = event.clipboardData?.getData('text/plain');
      if (!text) return false;

      if (openRef.current) {
        if (!text.includes(CAPSULE_BOUNDARY)) return false;

        event.preventDefault();
        const cleanText = text.replace(new RegExp(CAPSULE_BOUNDARY, 'g'), '');
        view.dispatch(view.state.replaceSelection(cleanText));
        return true;
      }

      event.preventDefault();
      const regex = new RegExp(AUTO_COMPLETE_CAPSULE_REGEX.source, 'g');
      const processedText = text.replace(regex, `$&${CAPSULE_BOUNDARY}`);
      view.dispatch(view.state.replaceSelection(processedText));
      return true;
    }, []);
```

已进行更改。

---

这段代码的作用是告诉 CodeMirror 编辑器：**由这个插件生成的装饰器（也就是那些胶囊）是“原子性”的（Atomic）。**

简单来说，它定义了光标和选区在遇到胶囊时的行为。

### 具体解释

1.  **`EditorView.atomicRanges`**：

    - 这是 CodeMirror 的一个特性（Facet），用于定义文档中哪些范围是**不可分割**的。
    - **原子范围（Atomic Range）** 意味着光标**不能停留在这个范围内部**。

2.  **`provide(plugin)`**：
    - 这是 `ViewPlugin` 的一种机制，允许插件向编辑器提供额外的功能。这里它将插件生成的 `decorations`（胶囊的位置信息）提供给了 `atomicRanges`。

### 这段代码带来的实际效果

如果没有这段代码，或者没有正确配置原子范围，用户体验会很糟糕：

1.  **光标移动**：

    - **有这段代码**：当你按左右方向键移动光标遇到胶囊时，光标会**直接跳过**整个胶囊，从左边直接跳到右边（或反之）。
    - **没有这段代码**：光标可能会“走进”胶囊内部（虽然视觉上胶囊是一个整体，但底层文本可能还是 `@Button`），导致光标在胶囊中间闪烁，或者位置计算错误。

2.  **删除行为**：

    - **有这段代码**：在胶囊旁边按 Backspace（退格键），会**一次性删除整个胶囊**。
    - **没有这段代码**：可能会只删除胶囊底层文本的一个字符（例如把 `@Button` 删成了 `@Butto`），导致胶囊突然失效变回普通文本，或者出现奇怪的渲染状态。

3.  **文本选择**：
    - 它确保你不能只选中胶囊名字的一半。你要么选中整个胶囊，要么不选中。

### 总结

这段代码是**“胶囊化”体验的核心保障**。它确保了胶囊在编辑器中表现得像一个**实体字符**（类似 Emoji 或图片），而不是一串可以被从中间切开的普通文本。

---

## forceUpdateEditor

在 CodeMirror 6 中，强制更新编辑器的最佳实践取决于你想要达到的目的。让我分析一下当前代码和推荐的做法：

## 当前代码分析

```tsx
const forceUpdateEditor = useCallback(() => {
  if (editorRef.current?.$view) {
    editorRef.current.$view.dispatch({})
  }
}, [])
```

这种方式是**可行的**，但不是最优的。`dispatch({})` 会触发一个空的状态更新，这会让编辑器重新渲染。

## CodeMirror 6 强制更新的最佳实践

### 1. **如果只是需要重绘视图（不改变状态）**

使用 `requestMeasure`：

```tsx
editorRef.current.$view.requestMeasure()
```

这是最轻量的方式，只会触发视图的重新测量和重绘。

### 2. **如果需要触发扩展/插件的更新**

使用带有空 effects 的 dispatch：

```tsx
editorRef.current.$view.dispatch({
  effects: StateEffect.appendConfig.of([]) // 空配置
})
```

或者使用 `StateEffect.reconfigure` 来触发更新。

### 3. **如果需要重新计算装饰（Decorations）**

```tsx
import { StateEffect } from '@codemirror/state'

// 创建一个自定义的 effect
const forceUpdate = StateEffect.define<null>()

// dispatch 时使用
editorRef.current.$view.dispatch({
  effects: forceUpdate.of(null)
})
```

### 4. **当前场景的建议**

对于你的代码，如果目的是强制刷新视图（比如外部状态变化后需要更新装饰），推荐：

```tsx
const forceUpdateEditor = useCallback(() => {
  const view = editorRef.current?.$view
  if (view) {
    // 方式一：最轻量，只重绘
    view.requestMeasure()

    // 方式二：如果需要触发状态更新（当前代码的方式也可以）
    // view.dispatch({});
  }
}, [])
```

## 总结

| 方法                         | 适用场景             | 性能          |
| ---------------------------- | -------------------- | ------------- |
| `view.requestMeasure()`      | 只需重绘视图         | ⭐⭐⭐ 最轻量 |
| `view.dispatch({})`          | 需要触发状态更新周期 | ⭐⭐ 轻量     |
| `dispatch({ effects: ... })` | 需要触发特定插件更新 | ⭐ 按需       |

你当前使用的 `dispatch({})` 是可接受的做法，如果只是为了触发重绘，可以考虑改用 `requestMeasure()` 更轻量。

这是因为 CodeMirror 的 `MatchDecorator` 机制默认只在文档内容发生变化时才更新装饰器（Decorations）。仅仅调用 `dispatch({})` 发送空事务并不会触发 `MatchDecorator` 重新扫描文档，因此 `MentionCapsuleWidget` 不会被重新创建，也就无法读取到最新的“已删除”状态。

我们需要定义一个自定义的 `StateEffect`，显式通知插件强制重新生成装饰器。

请按以下步骤修改：

### 1. 修改 capsule-plugin.ts

定义一个 `forceUpdateCapsulesEffect`，并在 `update` 方法中监听它。当接收到该 Effect 时，强制调用 `matcher.createDeco` 重新生成所有胶囊。

```typescript
import {
  Decoration,
  DecorationSet,
  EditorView,
  MatchDecorator,
  ViewPlugin,
  ViewUpdate
} from '@codemirror/view';
import { StateEffect } from '@codemirror/state';

import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget';
import type { CapsuleGlobalState, Extension } from '../types';
import { CAPSULE_REGEX, extractAllCapsules } from '../utils';
import type { ICapsuleInfo } from '@/components/ChatProvider';

const getCapsulesKey = (capsules: ICapsuleInfo[]): string => JSON.stringify(capsules.map((c) => c.id));

// 定义强制更新的 Effect
export const forceUpdateCapsulesEffect = StateEffect.define<null>();

export const createCapsulePlugin = (globalState: CapsuleGlobalState): Extension[] => {
  const matcher = new MatchDecorator({
    regexp: CAPSULE_REGEX,
    decoration: (match) => {
      const name = match[1];
      if (!name || name.length === 0) return null;
      return Decoration.replace({
        widget: new MentionCapsuleWidget(name, globalState),
        inclusive: false
      });
    }
  });

  const decorationPlugin = ViewPlugin.fromClass(
    class CapsuleDecorationPlugin {
      decorations: DecorationSet;

      constructor(view: EditorView) {
        this.decorations = matcher.createDeco(view);
      }

      update(update: ViewUpdate) {
        // 检查是否有强制更新的 Effect
        const forceUpdate = update.transactions.some((tr) =>
          tr.effects.some((e) => e.is(forceUpdateCapsulesEffect))
        );

        if (forceUpdate) {
          // 强制重新创建所有装饰器，这将触发 Widget 的重新实例化和状态检查
          this.decorations = matcher.createDeco(update.view);
        } else {
          this.decorations = matcher.updateDeco(update, this.decorations);
        }
      }
    },
    {
      decorations: (plugin) => plugin.decorations,
      provide(plugin) {
        return EditorView.atomicRanges.of((view) => view.plugin(plugin)?.decorations ?? Decoration.none);
      }
    }
  );

  let lastCapsulesKey = '';
  // ...existing code...
  const updateListener = EditorView.updateListener.of((update) => {
// ...existing code...
```

### 2. 修改 index.tsx

引入刚才定义的 `forceUpdateCapsulesEffect`，并在 `forceUpdateEditor` 中 dispatch 这个 Effect。

```typescript
// ...existing code...
import type { ITextAreaProps, CapsuleGlobalState } from './types';
import { extractMentionQuery, CAPSULE_BOUNDARY, AUTO_COMPLETE_CAPSULE_REGEX } from './utils';
import { createCapsulePlugin, createHistoryNavigationKeymap, editorTheme } from './extensions';
// 引入 Effect (如果 extensions/index.ts 没有导出，请直接从文件引入)
import { forceUpdateCapsulesEffect } from './extensions/capsule-plugin';
import { useHistoryNavigation, usePopupPosition } from './hooks';

const plugins = [...preset, domEventHandler('keydown')];

const Container = forwardRef<unknown, ITextAreaProps>(
// ...existing code...
    const closePanel = useCallback(() => {
      setOpen(false);
      openRef.current = false;
    }, []);

    const forceUpdateEditor = useCallback(() => {
      if (editorRef.current?.$view) {
        // 发送带有强制更新 Effect 的事务
        editorRef.current.$view.dispatch({
          effects: forceUpdateCapsulesEffect.of(null)
        });
      }
    }, []);
    useEffect(() => {
      emitter?.on('ecom_chat:force_update_editor', forceUpdateEditor);
// ...existing code...
```

**原理说明：**

1.  **`forceUpdateCapsulesEffect`**: 这是一个信号，告诉 CodeMirror 发生了一些与文档内容无关但影响显示的变化。
2.  **`matcher.createDeco(update.view)`**: 在插件的 `update` 方法中，当检测到这个信号时，我们不再使用增量更新 (`updateDeco`)，而是强制全量重新扫描 (`createDeco`)。
3.  **重新实例化 Widget**: 全量扫描会重新执行 `decoration` 回调，创建新的 `MentionCapsuleWidget` 实例。
4.  **状态更新**: 新实例在构造函数中调用 `getCommandState()`，此时能获取到最新的“已删除”状态（`isDisabled: true`）。
5.  **DOM 更新**: CodeMirror 比较新旧 Widget，发现状态不同（`eq` 返回 false），从而更新 DOM，显示禁用样式。

---

你的感觉很敏锐。目前的实现确实有点“分散”，因为 `forceUpdate` 的逻辑被拆分到了两个不同的生命周期中：

1.  **视图层（ViewPlugin）**：负责重绘 UI（比如变灰、变红）。
2.  **数据层（UpdateListener）**：负责通知外部组件（React）胶囊状态变了。

虽然它们各司其职，但为了让代码更内聚、更易维护，我们可以通过**重构**将核心逻辑集中起来。

### 为什么会分散？

CodeMirror 6 的架构决定了：

- **Decorations (UI)** 必须在 `ViewPlugin` 或 `StateField` 中生成。
- **Side Effects (回调)** 最好在 `EditorView.updateListener` 中处理。

### 优化方案：集中逻辑

我们可以把“计算胶囊数据”的逻辑提取出来，作为一个单一的真实数据源（Source of Truth），然后让 UI 和回调都依赖这个数据源。

#### 方案：使用 `StateField` 管理胶囊状态

创建一个 `StateField` 来专门存储当前的胶囊列表。这样，无论是 UI 渲染还是外部回调，都只需要监听这个 Field 的变化。

**重构后的 capsule-plugin.ts：**

```typescript
import {
  Decoration,
  DecorationSet,
  EditorView,
  MatchDecorator,
  ViewPlugin,
  ViewUpdate
} from '@codemirror/view'
import { StateEffect, StateField, Extension } from '@codemirror/state'
import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget'
import type { CapsuleGlobalState } from '../types'
import { CAPSULE_REGEX, extractAllCapsules } from '../utils'

export const forceUpdateCapsulesEffect = StateEffect.define<null>()

// 1. 定义一个 StateField 来集中管理胶囊数据
// 这是唯一的“真实数据源”
const capsuleStateField = (globalState: CapsuleGlobalState) =>
  StateField.define<DecorationSet>({
    create(state) {
      // 初始创建装饰
      return createDecorations(state.doc.toString(), globalState)
    },
    update(decorations, tr) {
      // 核心逻辑集中在这里：
      // 无论是文档变了，还是收到了强制更新信号，都重新计算装饰
      const forceUpdate = tr.effects.some(e => e.is(forceUpdateCapsulesEffect))

      if (tr.docChanged || forceUpdate) {
        return createDecorations(tr.state.doc.toString(), globalState)
      }

      // 必须 map 装饰以适应文档变化（即使 docChanged 为 false，selection 变化也可能触发 update）
      return decorations.map(tr.changes)
    },
    provide: f => EditorView.decorations.from(f) // 自动提供给视图渲染
  })

// 辅助函数：生成装饰集
function createDecorations(text: string, globalState: CapsuleGlobalState) {
  const widgets = []
  const regex = new RegExp(CAPSULE_REGEX, 'g')
  let match

  while ((match = regex.exec(text)) !== null) {
    const name = match[1]
    const start = match.index
    const end = start + match[0].length

    if (name) {
      widgets.push(
        Decoration.replace({
          widget: new MentionCapsuleWidget(name, globalState),
          inclusive: false
        }).range(start, end)
      )
    }
  }

  return Decoration.set(widgets)
}

export const createCapsulePlugin = (globalState: CapsuleGlobalState): Extension[] => {
  // 2. 监听器只负责“通知外部”，不再负责计算逻辑
  // 它只需要对比 StateField 的变化
  let lastCapsulesKey = ''

  const updateListener = EditorView.updateListener.of(update => {
    const forceUpdate = update.transactions.some(tr =>
      tr.effects.some(e => e.is(forceUpdateCapsulesEffect))
    )

    if (!update.docChanged && !forceUpdate) return

    // 重新计算胶囊数据用于回调（这里无法复用 DecorationSet，因为格式不同，但逻辑是同步的）
    const text = update.state.doc.toString()
    const allCommands = globalState.commandsGetter?.() || []
    const capsules = extractAllCapsules(text, allCommands)

    // 生成 Key (包含 disabled 状态)
    const capsulesKey = JSON.stringify(capsules.map(c => ({ id: c.id, disabled: c.disabled })))

    if (capsulesKey !== lastCapsulesKey) {
      lastCapsulesKey = capsulesKey
      // 放入宏任务避免 React 更新冲突
      setTimeout(() => {
        globalState.capsulesChangeHandler?.(capsules)
      })
    }
  })

  return [
    capsuleStateField(globalState), // 负责 UI 和状态管理
    updateListener // 负责副作用（回调）
  ]
}
```

### 这种写法的优点

1.  **逻辑归一**：`StateField` 的 `update` 方法成为了判断“是否需要更新胶囊”的唯一入口。无论是文档变动还是 `forceUpdate` 信号，都在这里处理。
2.  **自动渲染**：通过 `provide: EditorView.decorations.from(f)`，CodeMirror 会自动使用这个 Field 里的数据来渲染 UI，不需要再写一个 `ViewPlugin`。
3.  **更符合 CM6 哲学**：CodeMirror 6 推荐使用 `StateField` 来管理派生状态（Derived State）。

### 总结

之前的写法是用 `MatchDecorator` (ViewPlugin) 做 UI，用 `updateListener` 做回调，两者是平行的，所以逻辑分散。

现在的写法是用 `StateField` 做核心驱动，UI 自动从 StateField 衍生，回调监听 State 变化。虽然代码量差不多，但数据流向更清晰了。

---

是的，其实 `MatchDecorator` 是 CodeMirror 官方推荐处理这种正则替换的标准方式，它内部做了很多增量更新的优化（只扫描变动的部分），性能通常更好。

既然你觉得原来的结构更顺手，那我们就**保留原来的 `MatchDecorator` 结构**，只需要把刚才讨论的两个核心修复（`disabled` 状态检查、`forceUpdate` 逻辑）加进去即可。

这是基于“原来的结构”优化后的最终版本：

```typescript
import {
  Decoration,
  DecorationSet,
  EditorView,
  MatchDecorator,
  ViewPlugin,
  ViewUpdate
} from '@codemirror/view'
import { StateEffect, Extension } from '@codemirror/state'
import { MentionCapsuleWidget } from '../widgets/MentionCapsuleWidget'
import type { CapsuleGlobalState } from '../types'
import { CAPSULE_REGEX, extractAllCapsules } from '../utils'
import type { ICapsuleInfo } from '@/components/ChatProvider'

// 修复 1: Key 生成逻辑包含 disabled 状态
const getCapsulesKey = (capsules: ICapsuleInfo[]): string =>
  JSON.stringify(capsules.map(c => ({ id: c.id, disabled: c.disabled })))

export const forceUpdateCapsulesEffect = StateEffect.define<null>()

export const createCapsulePlugin = (globalState: CapsuleGlobalState): Extension[] => {
  // UI 部分：使用 MatchDecorator (性能更好)
  const matcher = new MatchDecorator({
    regexp: CAPSULE_REGEX,
    decoration: match => {
      const name = match[1]
      if (!name || name.length === 0) return null
      return Decoration.replace({
        widget: new MentionCapsuleWidget(name, globalState),
        inclusive: false
      })
    }
  })

  const decorationPlugin = ViewPlugin.fromClass(
    class CapsuleDecorationPlugin {
      decorations: DecorationSet

      constructor(view: EditorView) {
        this.decorations = matcher.createDeco(view)
      }

      update(update: ViewUpdate) {
        const forceUpdate = update.transactions.some(tr =>
          tr.effects.some(e => e.is(forceUpdateCapsulesEffect))
        )

        // 修复 2: 强制更新时，重新创建所有装饰以刷新 Widget 状态（如变灰）
        if (forceUpdate) {
          this.decorations = matcher.createDeco(update.view)
        } else {
          this.decorations = matcher.updateDeco(update, this.decorations)
        }
      }
    },
    {
      decorations: plugin => plugin.decorations,
      provide(plugin) {
        return EditorView.atomicRanges.of(
          view => view.plugin(plugin)?.decorations ?? Decoration.none
        )
      }
    }
  )

  // 逻辑部分：监听变化通知外部
  let lastCapsulesKey = ''
  const updateListener = EditorView.updateListener.of(update => {
    const forceUpdate = update.transactions.some(tr =>
      tr.effects.some(e => e.is(forceUpdateCapsulesEffect))
    )

    // 修复 3: 监听 docChanged 或 forceUpdate
    if (!update.docChanged && !forceUpdate) return

    const text = update.state.doc.toString()
    const allCommands = globalState.commandsGetter?.() || []
    const capsules = extractAllCapsules(text, allCommands)

    // 使用修复后的 Key 生成函数
    const capsulesKey = getCapsulesKey(capsules)

    if (capsulesKey !== lastCapsulesKey) {
      lastCapsulesKey = capsulesKey
      setTimeout(() => {
        globalState.capsulesChangeHandler?.(capsules)
      })
    }
  })

  return [decorationPlugin, updateListener]
}
```

### 这个版本的特点：

1.  **UI 渲染**：回归 `MatchDecorator`，利用 CM6 原生优化，处理大文档性能更好。
2.  **强制更新**：在 `ViewPlugin` 中处理 `forceUpdate`，确保 Widget 样式（变灰/变红）能刷新。
3.  **回调通知**：在 `updateListener` 中处理 `forceUpdate`，确保 React 层能收到状态变更通知。

---

这段代码的作用是：**当自动补全列表中的“当前选中项”发生变化时，强制刷新 CodeMirror 编辑器视图，以便更新界面上的“幽灵文本”或装饰器。**

让我们逐行拆解它的技术含义：

### 1. 监听变化 (Reactive Tracker)

```typescript
;() => {
  return [this.currentSelectedItem.value?.idName]
}
```

- 这是一个 MobX 风格的响应式追踪函数。
- 它监听 `this.currentSelectedItem`（当前选中的候选项）的 `idName`。
- **含义**：每当用户在自动补全菜单中按“上/下”键切换选项时，这个值就会变，从而触发后面的逻辑。

### 2. 执行副作用 (Effect)

```typescript
;() => {
  setTimeout(() => {
    this.editorView.dispatch({
      effects: []
    })
  }, 0)
}
```

- **`this.editorView.dispatch({ effects: [] })`**:

  - 这是 CodeMirror 6 中的一个常用技巧。它发送了一个**空事务**（没有修改文档内容，也没有移动光标）。
  - **目的**：强制触发 CodeMirror 的 **Update Cycle（视图更新周期）**。
  - CodeMirror 是数据驱动的，通常只有当 `EditorState` 变化时才会重绘。但这里的 `currentSelectedItem` 是外部状态（MobX 状态），CodeMirror 默认感知不到它的变化。
  - 通过手动 dispatch，通知编辑器重新计算所有的插件（Plugins）和装饰器（Decorations）。

- **`setTimeout(..., 0)`**:
  - 将 dispatch 操作推迟到下一个事件循环（Next Tick）。
  - **目的**：防止在当前的更新周期中再次触发更新，避免出现 "Apply view update recursively"（递归更新视图）的报错。

### 3. 为什么要这么做？

结合文件下方的代码：

```typescript
// ...
decorationRange = this.showCandidateText();
// ...
showCandidateText () {
   // ... 根据 currentSelectedItem 生成灰色的预览文字 ...
}
```

这个类似乎实现了一个功能：**在光标后面显示当前选中项的灰色预览（Ghost Text）**。

- **流程**：
  1.  用户按向下键 -> `currentSelectedItem` 变了。
  2.  MobX 状态更新了，但 CodeMirror 界面还没动，灰色的预览文字还是旧的。
  3.  `flushCurrentSelectedItemAction` 触发 -> 发送空 dispatch。
  4.  CodeMirror 收到 dispatch -> 触发重绘 -> 重新调用 `showCandidateText`。
  5.  界面上的灰色预览文字更新为最新选中的项。

**总结：这是一个连接“外部响应式状态（MobX）”与“编辑器内部渲染循环（CodeMirror）”的桥梁代码。**
