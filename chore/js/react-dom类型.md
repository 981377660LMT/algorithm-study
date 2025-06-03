# ReactDOM 类型详解

ReactDOM 是 React 的核心包之一，负责将 React 组件渲染到 DOM 环境中。下面我将详细讲解 ReactDOM 的类型定义。

## 命名空间和导入

```typescript
export as namespace ReactDOM

import {
  ReactInstance,
  Component,
  ComponentState,
  ReactElement,
  FunctionComponentElement,
  CElement,
  DOMAttributes,
  DOMElement,
  ReactNode,
  ReactPortal
} from 'react'
```

这表明 ReactDOM 被导出为一个命名空间，并且它依赖于 React 包中的多种类型。

## 核心 API

### 渲染相关

```typescript
export const version: string
export const render: Renderer
export const hydrate: Renderer
```

- `version`: ReactDOM 的版本号
- `render`: 核心渲染函数，将 React 元素渲染到指定 DOM 容器
- `hydrate`: 用于服务端渲染后的客户端激活（水合）过程

### DOM 操作

```typescript
export function findDOMNode(instance: ReactInstance | null | undefined): Element | null | Text
export function unmountComponentAtNode(container: Element | DocumentFragment): boolean
```

- `findDOMNode`: 查找组件实例对应的 DOM 节点（已废弃）
- `unmountComponentAtNode`: 从 DOM 中卸载组件

### Portal 相关

```typescript
export function createPortal(
  children: ReactNode,
  container: Element | DocumentFragment,
  key?: null | string
): ReactPortal
```

- `createPortal`: 创建 Portal，允许将子元素渲染到父组件 DOM 层次结构之外

## Renderer 类型详解

`Renderer` 是一个复杂的接口类型，定义了渲染函数的各种重载形式：

```typescript
export interface Renderer {
  <T extends Element>(
    element: DOMElement<DOMAttributes<T>, T>,
    container: Container | null,
    callback?: () => void
  ): T

  // ... 其他重载
}
```

这个接口有多个重载，分别处理：

- DOM 元素渲染
- 函数组件元素渲染
- 类组件元素渲染
- 元素数组渲染

返回类型根据渲染内容的不同而变化，可能是 DOM 元素、组件实例或 void。

## 容器类型

```typescript
export type Container = Element | Document | DocumentFragment
```

`Container` 类型定义了可作为渲染目标的 DOM 容器类型，包括普通 DOM 元素、文档对象或文档片段。

## 批处理与同步更新

```typescript
export function flushSync<R>(fn: () => R): R
export function flushSync<A, R>(fn: (a: A) => R, a: A): R

export function unstable_batchedUpdates<A, R>(callback: (a: A) => R, a: A): R
export function unstable_batchedUpdates<R>(callback: () => R): R
```

- `flushSync`: 强制同步执行更新并立即刷新，绕过批处理机制
- `unstable_batchedUpdates`: 将多个更新批量处理，减少渲染次数

## 不稳定 API

```typescript
export function unstable_renderSubtreeIntoContainer<T extends Element>(
  parentComponent: Component<any>,
  element: DOMElement<DOMAttributes<T>, T>,
  container: Element,
  callback?: (element: T) => any
): T
// ... 其他重载
```

`unstable_renderSubtreeIntoContainer` 是一个不稳定 API，用于在特定父组件上下文中渲染子树到指定容器。有多个重载以支持不同类型的元素和组件。

## 类型特点

1. **泛型广泛使用**：为了提供类型安全，同时保持 API 的灵活性
2. **函数重载**：多个重载签名处理不同参数组合和返回类型

3. **条件类型**：根据输入参数类型确定返回类型

4. **联合类型**：如 `Container = Element | Document | DocumentFragment`，增加 API 灵活性

5. **可选参数**：如回调函数通常是可选的 `callback?: () => void`

## React 18 新增 API

React 18 中，`createRoot` 和 `hydrateRoot` API 成为推荐的渲染方式，取代传统的 `render` 和 `hydrate`，但在你提供的类型定义中没有包含，可能是较旧版本的定义。

## 注意事项

- 标记为 `unstable_` 的 API 不保证稳定性，可能在未来版本中变动
- `findDOMNode` 已被官方废弃，推荐使用 ref 代替
- React 18 中，批处理行为已发生变化，默认启用自动批处理

ReactDOM 的类型定义反映了它作为连接 React 与实际 DOM 环境的桥梁角色，类型设计精确捕获了各种 API 的行为和约束。
