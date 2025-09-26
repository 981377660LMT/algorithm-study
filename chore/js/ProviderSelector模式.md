# `分层Hook`

关键是：灵活拆分 hook，重计算的可以抽到父组件

```ts
const relationMap = useRelationMap() // 父组件
const relation = useRelation(relationMap, name) // 子组件
```

---

计算和缓存提升到所有组件共享的更高层级。一个好的做法是创建一个新的全局状态或者一个专门的 hook
这种技巧的本质是：**将昂贵的、可共享的计算从消费它的多个组件中分离出来，提升到更高层级进行统一计算和缓存，然后将计算结果分发给消费者。**
这是一种典型的性能优化模式，可以看作是**“计算的关注点分离”**（Separation of Concerns for Computation）。

### 问题根源

在原始代码中，昂贵的计算（遍历并构建 `relationMap`）与消费该计算结果的逻辑（从 `Map` 中获取特定 `name` 的数据）耦合在同一个 `useCodeRelations` Hook 中。当多个组件调用这个 Hook 时，每个组件都会独立地、重复地执行这个昂贵的计算，造成了资源浪费。

### 通用模式抽象

我们可以将这个模式抽象为 **“Provider/Selector” 模式** 或 **“分层 Hook” 模式**。

这个模式包含两个核心部分：

1.  **数据提供者 (Provider Hook / Centralized Computation Hook)**

    - **职责**: 负责执行昂贵的、一次性的计算或数据转换。
    - **实现**:
      - 它通常是一个自定义 Hook（例如我们创建的 `useAllCodeRelations`）。
      - 它会订阅一个原始数据源（如全局状态、Context、或 API 响应）。
      - 它使用 `useMemo` 来缓存计算结果。只有当原始数据源发生变化时，它才会重新计算。
    - **目标**: 确保在整个应用或某个子树中，无论有多少消费者，昂贵的计算只执行一次。

2.  **数据消费者 (Selector Hook / Consumer Component)**
    - **职责**: 从“数据提供者”获取已经计算好的数据，并根据自身需求进行轻量级的筛选或派生。
    - **实现**:
      - 它可以是另一个自定义 Hook（如 `useCodeRelations`）或一个组件。
      - 它调用“数据提供者” Hook 来获取共享的、已缓存的数据。
      - 它可能会使用另一个 `useMemo` 来缓存从共享数据中派生出的自身所需数据，但这步操作通常非常快（例如，从 Map 中 `get` 一个值）。
    - **目标**: 快速、高效地获取所需数据，而不触发任何昂贵的重复计算。

### 如何应用这个通用模式

当你发现以下情况时，就可以考虑使用此模式：

1.  **识别昂贵计算**: 你有一个函数或一段逻辑，它需要遍历大数组、处理复杂对象或执行其他耗时操作。
2.  **识别共享数据源**: 这个昂贵计算的输入是来自一个所有组件都能访问的共享数据源（例如 Redux/Zustand store、React Context、全局单例）。
3.  **识别多个消费者**: 应用中有多个组件或 Hook 都需要基于这个计算结果来展示信息或执行逻辑。

**抽象步骤：**

1.  **创建 `useSharedData` Hook (Provider)**:

    - 将昂贵的计算逻辑从你的组件或现有 Hook 中抽离出来。
    - 创建一个新的 Hook，例如 `useProcessedSharedData`。
    - 在这个新 Hook 内部，获取原始数据，并用 `useMemo` 包裹你的昂贵计算逻辑。
    - 返回 `useMemo` 的结果。

2.  **重构消费者**:
    - 在原来的组件或 Hook 中，删除昂贵的计算逻辑。
    - 转而调用你新创建的 `useProcessedSharedData` Hook 来直接获取计算好的数据。
    - 在此基础上进行后续的、轻量级的操作。

这个模式不仅限于 Hook，它也是 Redux 生态中 `reselect` 库的核心思想，即创建可组合、可记忆化的选择器（Selectors），以避免在状态树变化时进行不必要的派生数据计算。通过自定义 Hook，我们在 React 中原生实现了类似的思想。

---

你指出的问题非常关键，完全正确。

你遇到的问题是，即使我们将逻辑拆分成了 `useRelationMap` 和 `useRelation`，但由于每个调用 `useRelation` 的组件都会独立调用 `useRelationMap`，这导致 `useRelationMap` Hook 本身被实例化多次。因此，它内部的 `useMemo` 也在每个组件实例中都独立运行了一次，昂贵的计算并没有被真正地共享。

这是 React Hook 的一个核心行为：**Hook 在每个调用它的组件中都有独立的状态和 Memoization。**

### 解决方案：使用 React Context

为了在所有组件之间**真正地共享一个**计算结果，我们需要将这个结果提升到一个**所有组件都能访问的共同祖先上**。实现这一点的标准 React 模式就是 **Context**。

我们将创建一个 `RelationProvider`，它负责执行一次昂贵的计算，然后通过 Context 将计算结果（`relationMap`）提供给其下的所有子组件。然后，`useRelation` Hook 就可以直接从这个 Context 中消费已经计算好的数据。

下面是具体的实现步骤：

#### 1. 创建 Context 和 Provider

我们需要创建一个新的文件来存放 Context 和 Provider，或者将它放在 `util.tsx` 的顶部。

```tsx
import React, { createContext, useContext, useMemo } from 'react';
// ... 其他 import

// 1. 创建一个 Context 来持有计算好的 relationMap
const RelationMapContext = createContext<Map<string, { quoteTo: Set<string>; quoteBy: Set<string> }> | null>(null);

// 2. 创建一个 Provider 组件
// 这个组件应该被包裹在所有需要使用 useRelation 的组件的共同祖先上
export const RelationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const relations = useGlobalStateRelations((a) => a);

  // 在 Provider 中执行一次昂贵的计算
  const relationMap = useMemo(() => {
    console.log('🔴 昂贵的 Relation Map 计算正在执行... (这应该只发生一次)');
    const map = new Map<string, { quoteTo: Set<string>; quoteBy: Set<string> }>();
    const relationList = Object.values(relations) as { v: string; w: string }[];
    for (const item of relationList) {
      const vName = item.v.split('.')[0];
      const wName = item.w.split('.')[0];
      if (vName === wName) {
        continue;
      }
      if (!map.has(vName)) {
        map.set(vName, { quoteTo: new Set(), quoteBy: new Set() });
      }
      map.get(vName)!.quoteBy.add(wName);
      if (!map.has(wName)) {
        map.set(wName, { quoteTo: new Set(), quoteBy: new Set() });
      }
      map.get(wName)!.quoteTo.add(vName);
    }
    return map;
  }, [relations]);

  return (
    <RelationMapContext.Provider value={relationMap}>
      {children}
    </RelationMapContext.Provider>
  );
};

// 3. 创建一个自定义 Hook 来方便地消费 Context
const useRelationMap = () => {
  const context = useContext(RelationMapContext);
  if (!context) {
    throw new Error('useRelationMap 必须在 RelationProvider 内部使用');
  }
  return context;
};

// 4. 修改 useRelation Hook，让它从 Context 中消费数据
export const useRelation = (name: string) => {
  const relationMap = useRelationMap(); // 现在这个 Hook 非常轻量，只是读取 Context

  const res = useMemo(() => {
    const relation = relationMap.get(name);
    return {
      quoteTo: relation ? Array.from(relation.quoteTo) : [],
      quoteBy: relation ? Array.from(relation.quoteBy) : [],
    };
  }, [relationMap, name]);

  return res;
};


// ... 文件中原有的其他代码 ...
/**
 * @deprecated Use {@link useRelation} instead
 */
export const getRelationsByName = ...
// ...
```

#### 2. 在应用中使用 `RelationProvider`

现在，你需要找到一个合适的祖先组件，用 `RelationProvider` 将所有可能调用 `useRelation` 的组件包裹起来。例如，如果这些组件都在一个调试面板里，你可以包裹整个面板。

```tsx
// 这是一个示例，你需要找到你项目中合适的包裹位置
// 比如在 JsonViewerWrapper 的某个父组件中

import { RelationProvider } from './path/to/your/relation/util'

const DebugPanel = () => {
  return (
    <RelationProvider>
      {/* 这里是你的应用的其他部分，
          所有在这里渲染的、调用了 useRelation 的子组件
          都将共享同一个 relationMap 实例 */}
      <YourComponentTree />
    </RelationProvider>
  )
}
```

通过这种方式，昂贵的 `relationMap` 计算只在 `RelationProvider` 组件中执行一次，其结果通过 Context 高效地分发给所有消费者，彻底解决了重复计算的问题。
