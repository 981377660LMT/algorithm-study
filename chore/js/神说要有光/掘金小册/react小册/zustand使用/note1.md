## 1. 基础概念

Zustand 是一个轻量级的 React 状态管理库，核心理念是简单、直观的状态管理。

## 2. 基本使用方法

### 创建 Store

```typescript
import { create } from 'zustand'

// 定义状态接口
interface BearState {
  bears: number
  increase: (by: number) => void
  decrease: (by: number) => void
  reset: () => void
}

// 创建 store
const useBearStore = create<BearState>((set, get) => ({
  bears: 0,
  increase: by => set(state => ({ bears: state.bears + by })),
  decrease: by => set(state => ({ bears: state.bears - by })),
  reset: () => set({ bears: 0 })
}))
```

### 在组件中使用

```typescript
// 使用整个状态
function BearCounter() {
  const { bears, increase, decrease, reset } = useBearStore()

  return (
    <div>
      <h1>{bears} around here...</h1>
      <button onClick={() => increase(1)}>one up</button>
      <button onClick={() => decrease(1)}>one down</button>
      <button onClick={reset}>reset</button>
    </div>
  )
}

// 使用选择器（推荐）
function BearCounter() {
  const bears = useBearStore(state => state.bears)
  const increase = useBearStore(state => state.increase)

  return (
    <div>
      <h1>{bears} around here...</h1>
      <button onClick={() => increase(1)}>one up</button>
    </div>
  )
}
```

## 3. 高级用法

### 异步操作

```typescript
interface UserState {
  users: User[]
  loading: boolean
  fetchUsers: () => Promise<void>
}

const useUserStore = create<UserState>((set, get) => ({
  users: [],
  loading: false,
  fetchUsers: async () => {
    set({ loading: true })
    try {
      const users = await api.getUsers()
      set({ users, loading: false })
    } catch (error) {
      set({ loading: false })
      console.error('Failed to fetch users:', error)
    }
  }
}))
```

### 嵌套状态更新

```typescript
interface NestedState {
  user: {
    name: string
    profile: {
      age: number
      email: string
    }
  }
  updateUserName: (name: string) => void
  updateUserAge: (age: number) => void
}

const useNestedStore = create<NestedState>(set => ({
  user: {
    name: '',
    profile: {
      age: 0,
      email: ''
    }
  },
  updateUserName: name =>
    set(state => ({
      user: {
        ...state.user,
        name
      }
    })),
  updateUserAge: age =>
    set(state => ({
      user: {
        ...state.user,
        profile: {
          ...state.user.profile,
          age
        }
      }
    }))
}))
```

### 使用 Immer 简化嵌套更新

```typescript
import { create } from 'zustand'
import { immer } from 'zustand/middleware/immer'

const useImmerStore = create<NestedState>()(
  immer(set => ({
    user: {
      name: '',
      profile: { age: 0, email: '' }
    },
    updateUserName: name =>
      set(state => {
        state.user.name = name
      }),
    updateUserAge: age =>
      set(state => {
        state.user.profile.age = age
      })
  }))
)
```

## 4. 中间件使用

### 持久化存储

```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

const usePersistStore = create<BearState>()(
  persist(
    set => ({
      bears: 0,
      increase: by => set(state => ({ bears: state.bears + by })),
      decrease: by => set(state => ({ bears: state.bears - by })),
      reset: () => set({ bears: 0 })
    }),
    {
      name: 'bear-storage', // localStorage key
      // 可选配置
      partialize: state => ({ bears: state.bears }) // 只持久化部分状态
    }
  )
)
```

### 开发工具集成

```typescript
import { create } from 'zustand'
import { devtools } from 'zustand/middleware'

const useDevStore = create<BearState>()(
  devtools(
    set => ({
      bears: 0,
      increase: by => set(state => ({ bears: state.bears + by }), false, 'increase'),
      decrease: by => set(state => ({ bears: state.bears - by }), false, 'decrease'),
      reset: () => set({ bears: 0 }, false, 'reset')
    }),
    {
      name: 'bear-store'
    }
  )
)
```

## 5. 类型定义解析

根据您提供的类型文件，理解关键类型：

```typescript
// ExtractState: 提取 Store 中的状态类型
type ExtractState<S> = S extends {
  getState: () => infer T
}
  ? T
  : never

// UseBoundStore: 绑定后的 Store 类型，支持多种调用方式
export type UseBoundStore<S> = {
  (): ExtractState<S> // 无参调用，返回完整状态
  <U>(selector: (state: ExtractState<S>) => U): U // 选择器调用
} & S
```

## 6. 最佳实践

### 状态分片

```typescript
// 将大型状态分割为小的 store
const useAuthStore = create<AuthState>(set => ({
  user: null,
  login: async credentials => {
    /* ... */
  },
  logout: () => set({ user: null })
}))

const useUIStore = create<UIState>(set => ({
  theme: 'light',
  sidebarOpen: false,
  toggleSidebar: () => set(state => ({ sidebarOpen: !state.sidebarOpen }))
}))
```

### 选择器优化

```typescript
// 避免不必要的重渲染
const BearCounter = () => {
  // ✅ 好的做法：只选择需要的状态
  const bears = useBearStore(state => state.bears)

  // ❌ 避免：选择整个状态对象
  // const { bears } = useBearStore()

  return <div>{bears}</div>
}

// 使用浅比较选择器
import { shallow } from 'zustand/shallow'

const BearActions = () => {
  const { increase, decrease } = useBearStore(
    state => ({ increase: state.increase, decrease: state.decrease }),
    shallow
  )

  return (
    <div>
      <button onClick={() => increase(1)}>+</button>
      <button onClick={() => decrease(1)}>-</button>
    </div>
  )
}
```

### 在组件外使用

```typescript
// 可以在组件外部直接访问 store
const store = useBearStore.getState()
console.log(store.bears) // 获取当前状态
store.increase(1) // 调用 action

// 订阅状态变化
const unsubscribe = useBearStore.subscribe(state => console.log('Bears changed:', state.bears))
```

## 7. 与其他状态管理库的对比

| 特性            | Zustand     | Redux  | Context API |
| --------------- | ----------- | ------ | ----------- |
| 包大小          | 很小 (~2kb) | 大     | 内置        |
| 学习曲线        | 平缓        | 陡峭   | 中等        |
| TypeScript 支持 | 优秀        | 优秀   | 良好        |
| 中间件生态      | 丰富        | 最丰富 | 有限        |
| 性能            | 优秀        | 优秀   | 需要优化    |

Zustand 的优势在于简单直观、类型安全、性能优秀，特别适合中小型项目的状态管理需求。
