# Vue.js 响应式原理详解

## 响应式系统整体架构

```mermaid
flowchart TD
    subgraph "数据层"
        rawData[原始数据对象] --> reactive[响应式代理]
    end

    subgraph "响应式系统核心"
        reactive <--> track[依赖追踪]
        reactive <--> trigger[更新触发]
        track --- effectsMap[依赖映射表]
        trigger --- effectsMap
    end

    subgraph "响应消费者"
        effectsMap --> render[组件渲染函数]
        effectsMap --> computed[计算属性]
        effectsMap --> watch[侦听器]
    end

    user[用户交互/数据修改] --> reactive
    render --> UI[用户界面]
```

## 依赖收集与更新触发流程

```mermaid
sequenceDiagram
    participant Component as 组件
    participant Render as 渲染函数
    participant Proxy as 响应式代理
    participant Dep as 依赖管理器

    Note over Component,Dep: 初始化 & 依赖收集阶段
    Component->>Render: 初始渲染
    activate Render
    Note right of Render: 将当前渲染函数设为activeEffect
    Render->>Proxy: 读取属性 state.count
    Proxy->>Dep: track(target, 'count')
    Dep->>Dep: 将activeEffect添加到依赖集合
    Proxy-->>Render: 返回属性值
    Render-->>Component: 完成渲染
    deactivate Render

    Note over Component,Dep: 更新触发阶段
    Component->>Proxy: 修改数据 state.count++
    Proxy->>Dep: trigger(target, 'count')
    Dep->>Render: 通知相关effect重新执行
    activate Render
    Render->>Component: 重新渲染组件
    deactivate Render
```

## 响应式系统类图

```mermaid
classDiagram
    class ReactiveEffect {
        +fn: Function
        +deps: Set[]
        +active: boolean
        +run()
        +stop()
    }

    class TargetMap {
        -WeakMap~target, Map~key, Set~ReactiveEffect~~
        +track(target, key)
        +trigger(target, key)
    }

    class ReactiveProxy {
        -target: Object
        -handlers: ProxyHandler
        +get(target, key, receiver)
        +set(target, key, value, receiver)
        +deleteProperty(target, key)
    }

    class ComputedRefImpl {
        -getter: Function
        -_value: any
        -_dirty: boolean
        +value
        +effect: ReactiveEffect
    }

    ReactiveEffect -- TargetMap : 注册到 >
    ReactiveProxy --> TargetMap : 调用track/trigger
    ComputedRefImpl --> ReactiveEffect : 使用
```

## 响应式系统详细实现解析

### 1. Proxy拦截与代理

响应式系统的核心是通过Proxy拦截数据的读取和修改操作：

- **get拦截器**：当数据被读取时，记录依赖关系（谁正在使用这个数据）
- **set拦截器**：当数据被修改时，通知所有依赖这个数据的消费者更新

### 2. 依赖收集与存储结构

```mermaid
flowchart LR
    subgraph "WeakMap - targetMap"
        target1[目标对象1] --> depsMap1[Map - depsMap]
        target2[目标对象2] --> depsMap2[Map - depsMap]
    end

    subgraph "Map - depsMap1"
        key1[属性key1] --> deps1[Set - deps]
        key2[属性key2] --> deps2[Set - deps]
    end

    subgraph "Set - deps1"
        effect1[effect函数1]
        effect2[effect函数2]
    end
```

### 3. 响应式工作原理细节

1. **初始化阶段**：

   - 通过`reactive()`创建数据的Proxy代理
   - 设置拦截器监听数据操作

2. **依赖收集阶段**：

   - 组件渲染时，effect函数执行
   - 访问响应式数据触发get拦截器
   - 通过track()建立数据与当前effect的依赖关系

3. **数据变化阶段**：

   - 修改响应式数据触发set拦截器
   - 通过trigger()查找依赖此数据的所有effect
   - 执行这些effect函数，更新UI

4. **清理阶段**：
   - 组件卸载时，清理相关effect防止内存泄漏
   - 数据未被使用时，依赖关系可被垃圾回收

通过以上机制，Vue实现了声明式、细粒度的自动UI更新，开发者只需关注数据变化，而无需手动操作DOM。
