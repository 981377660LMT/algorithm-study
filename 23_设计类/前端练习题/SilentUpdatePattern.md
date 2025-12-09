## 🎯 问题本质：循环触发链

### 没有 `isInternalUpdate` 时的死循环

```
1. 父组件传入新的 codeDetail (比如从 API 加载)
   ↓
2. useEffect 检测到 codeDetail 变化
   ↓
3. 执行 formInstance.setValues(initialValues)
   ↓
4. ⚠️ Formily 内部触发 onFormValuesChange 事件
   ↓
5. 你的 effects 回调执行：
   - 读取表单值
   - 调用 onSave(codeDetail.id, result)
   ↓
6. 父组件接收到 onSave 回调
   ↓
7. 父组件可能更新状态 (比如 setState)
   ↓
8. 父组件重新渲染，传入新的 codeDetail 引用
   ↓
9. 回到步骤 1，无限循环 🔁
```

### 关键点：**表单库的事件机制**

```typescript
formInstance.setValues(initialValues)
```

大多数表单库（包括 Formily、react-hook-form、antd Form）的 `setValues` 方法会：

1. **更新内部状态**
2. **触发 `onChange` / `onValuesChange` 等事件** ← 关键！

这是合理的设计，因为：

- 表单库不知道你是"手动 API 回填"还是"用户输入"
- 它需要保证数据流的一致性（所有订阅者都能感知变化）
- 它需要触发字段级联动、校验等逻辑

---

## 🔒 加锁后的流程

### 有 `isInternalUpdate` 时的正常流程

```
1. 父组件传入新的 codeDetail
   ↓
2. useEffect 检测到变化
   ↓
3. 设置 isInternalUpdate.current = true ← 🔑 加锁
   ↓
4. 执行 formInstance.setValues(initialValues)
   ↓
5. Formily 触发 onFormValuesChange
   ↓
6. effects 回调执行：
   if (isInternalUpdate.current) {
     return; // ← 🛑 提前退出，不调用 onSave
   }
   ↓
7. setTimeout 0ms 后解锁
   ↓
8. 流程结束，没有触发 onSave，没有循环 ✅
```

---

## 📊 对比分析

### 场景 A：用户手动输入

```typescript
// 用户在输入框输入 "hello"
<input onChange={e => form.setFieldValue('name', e.target.value)} />
```

**流程**：

```
用户输入
  ↓
formInstance 内部状态更新
  ↓
触发 onFormValuesChange
  ↓
此时 isInternalUpdate.current === false (没有被锁住)
  ↓
执行 onSave，保存到服务器 ✅
```

### 场景 B：外部数据回填

```typescript
// useEffect 中同步外部数据
formInstance.setValues(initialValues)
```

**流程**：

```
useEffect 执行
  ↓
设置 isInternalUpdate.current = true
  ↓
formInstance 内部状态更新
  ↓
触发 onFormValuesChange
  ↓
此时 isInternalUpdate.current === true (被锁住)
  ↓
提前 return，不执行 onSave ✅
```

---

## 🧠 本质原因总结

### 1. **表单库的"不可知性"**

表单库无法区分数据变更的来源是：

- 用户交互（键盘输入、点击、拖拽）
- 程序代码（API 回填、计算字段、格式化）

它只能一视同仁地触发 `onChange` 事件。

### 2. **React 的单向数据流悖论**

React 推崇单向数据流：

```
Props (codeDetail) → State (formValues) → UI
```

但在表单场景中，我们需要"双向同步"：

```
外部数据 ⇄ 表单状态
```

这会导致：

- 外部数据变化 → 更新表单 → 触发 onChange → 又更新外部数据 → 循环

### 3. **副作用的不可控性**

你的 `onSave` 回调可能触发父组件的任何逻辑：

```typescript
// 父组件
const handleSave = (id, values) => {
  // 可能的操作：
  setCodeDetail(values) // 触发重新渲染
  await api.save(values) // 异步操作
  showNotification() // 其他副作用
}
```

如果不加锁，每次外部数据同步都会触发这些逻辑，造成：

- **性能问题**：重复的 API 请求
- **逻辑错误**：无限循环、状态不一致
- **用户体验**：闪烁、卡顿

---

## 🏆 最佳实践提炼

### 核心原则：**区分数据变更的意图**

```typescript
// ✅ 正确模式
const isUserAction = useRef(false)

// 用户操作入口
const handleUserChange = value => {
  isUserAction.current = true
  formInstance.setFieldValue('xxx', value)
  isUserAction.current = false
}

// 系统同步入口
const syncFromExternal = values => {
  isUserAction.current = false // 显式标记
  formInstance.setValues(values)
}

// 统一的监听器
onFormValuesChange(() => {
  if (!isUserAction.current) {
    return // 忽略非用户操作
  }
  triggerSideEffects()
})
```

### 为什么 `setTimeout(0)` 能工作？

```typescript
isInternalUpdate.current = true
formInstance.setValues(initialValues) // 同步执行

// setValues 内部：
//   1. 更新状态
//   2. 同步触发 onFormValuesChange
//   3. 你的回调检查到 isInternalUpdate === true，提前返回

setTimeout(() => {
  isInternalUpdate.current = false // 在下一个事件循环解锁
}, 0)
```

**关键**：

- `setValues` 是**同步**触发 `onFormValuesChange` 的
- 在 `setTimeout` 的回调执行前，`onFormValuesChange` 已经完成
- 所以锁在正确的时机生效和释放

---

## 🚨 你当前代码的潜在风险

虽然加了 `isInternalUpdate` 解决了循环，但仍有隐患：

### 风险 1：时序竞态

```typescript
setTimeout(() => {
  isInternalUpdate.current = false
}, 0)
```

如果用户在这 0ms 内有输入（虽然几乎不可能，但高频操作时有概率），可能导致：

- 用户输入被误判为"内部更新"
- 用户输入没有触发 `onSave`

### 风险 2：依赖过度触发

```typescript
useEffect(() => {
  // ...
}, [codeDetail, initialValues, formInstance])
```

每次 `codeDetail` 引用变化（即使内容相同）都会触发，可能导致：

- 不必要的 `setValues` 调用
- 表单失去焦点或光标跳动

### 建议改进

```typescript
useEffect(() => {
  const formValues = formInstance.getValuesIn('')
  if (!isEqual(formValues, initialValues)) {
    isInternalUpdate.current = true
    formInstance.setValues(initialValues)
    // 改用 queueMicrotask 或在 setValues 的 Promise 后解锁
    queueMicrotask(() => {
      isInternalUpdate.current = false
    })
  }
}, [codeDetail, initialValues, formInstance])
```

或者更彻底的方案：

```typescript
// 记录上次的 codeDetail，只在真正变化时才同步
const prevCodeDetailRef = useRef(codeDetail)

useEffect(() => {
  if (isEqual(prevCodeDetailRef.current, codeDetail)) {
    return // 内容没变，不需要同步
  }

  prevCodeDetailRef.current = codeDetail

  isInternalUpdate.current = true
  formInstance.setValues(transferSchemaToForm(codeDetail))
  queueMicrotask(() => {
    isInternalUpdate.current = false
  })
}, [codeDetail, formInstance])
```

---

# 🎯 表单双向同步的通用 TypeScript 模式抽象

## 模式 1：互斥锁包装器（Mutex Wrapper）

适用于任何需要区分"用户操作"与"系统同步"的场景。

```typescript
/**
 * 互斥锁包装器 - 用于防止循环更新
 * @example
 * const mutex = new MutexWrapper();
 *
 * // 系统同步时加锁
 * mutex.withLock(() => {
 *   formInstance.setValues(data);
 * });
 *
 * // 监听器中检查
 * onFormChange(() => {
 *   if (mutex.isLocked) return;
 *   saveToServer();
 * });
 */
export class MutexWrapper {
  private lockCount = 0

  /** 当前是否处于锁定状态 */
  get isLocked(): boolean {
    return this.lockCount > 0
  }

  /** 执行加锁操作，自动在函数执行后解锁 */
  withLock<T>(fn: () => T): T {
    this.lockCount++
    try {
      return fn()
    } finally {
      // 使用 queueMicrotask 确保在当前同步流程后解锁
      queueMicrotask(() => {
        this.lockCount = Math.max(0, this.lockCount - 1)
      })
    }
  }

  /** 异步版本 */
  async withLockAsync<T>(fn: () => Promise<T>): Promise<T> {
    this.lockCount++
    try {
      return await fn()
    } finally {
      queueMicrotask(() => {
        this.lockCount = Math.max(0, this.lockCount - 1)
      })
    }
  }

  /** 手动加锁（需要配合 unlock 使用） */
  lock(): void {
    this.lockCount++
  }

  /** 手动解锁 */
  unlock(): void {
    queueMicrotask(() => {
      this.lockCount = Math.max(0, this.lockCount - 1)
    })
  }

  /** 重置锁状态 */
  reset(): void {
    this.lockCount = 0
  }
}
```

---

## 模式 2：React Hook - `useMutexState`

封装带互斥锁的状态管理。

```typescript
import { useState, useRef, useCallback } from 'react'

export interface MutexStateReturn<T> {
  /** 当前值 */
  value: T
  /** 用户主动变更（触发副作用） */
  setUserValue: (value: T | ((prev: T) => T)) => void
  /** 系统静默更新（不触发副作用） */
  setSystemValue: (value: T | ((prev: T) => T)) => void
  /** 当前是否为系统更新 */
  isSystemUpdate: () => boolean
}

/**
 * 带互斥锁的状态管理 Hook
 * @example
 * const { value, setUserValue, setSystemValue, isSystemUpdate } = useMutexState('');
 *
 * useEffect(() => {
 *   if (isSystemUpdate()) return; // 跳过系统更新
 *   saveToAPI(value);
 * }, [value]);
 */
export function useMutexState<T>(initialValue: T): MutexStateReturn<T> {
  const [value, setValue] = useState<T>(initialValue)
  const isSystemRef = useRef(false)

  const setUserValue = useCallback((newValue: T | ((prev: T) => T)) => {
    isSystemRef.current = false
    setValue(newValue)
  }, [])

  const setSystemValue = useCallback((newValue: T | ((prev: T) => T)) => {
    isSystemRef.current = true
    setValue(newValue)
    // 在下一个微任务中重置标志
    queueMicrotask(() => {
      isSystemRef.current = false
    })
  }, [])

  const isSystemUpdate = useCallback(() => {
    return isSystemRef.current
  }, [])

  return {
    value,
    setUserValue,
    setSystemValue,
    isSystemUpdate
  }
}
```

---

## 模式 3：表单同步控制器（Form Sync Controller）

专门用于处理表单的外部数据同步。

```typescript
import { isEqual } from 'lodash-es'

export interface FormSyncOptions<TExternal, TForm> {
  /** 外部数据转表单数据 */
  externalToForm: (external: TExternal) => TForm
  /** 表单数据转外部数据 */
  formToExternal: (form: TForm, originalExternal: TExternal) => Partial<TExternal>
  /** 表单实例的 setValues 方法 */
  setFormValues: (values: TForm) => void
  /** 表单实例的 getValues 方法 */
  getFormValues: () => TForm
  /** 保存回调 */
  onSave: (values: Partial<TExternal>) => void
}

/**
 * 表单同步控制器
 * 解决表单与外部数据双向绑定的循环问题
 */
export class FormSyncController<TExternal, TForm> {
  private mutex = new MutexWrapper()
  private externalDataRef: TExternal
  private lastSavedFormValuesRef: TForm | null = null

  constructor(initialExternal: TExternal, private options: FormSyncOptions<TExternal, TForm>) {
    this.externalDataRef = initialExternal
  }

  /**
   * 同步外部数据到表单
   * @param newExternal 新的外部数据
   * @param force 是否强制同步（忽略相等性检查）
   */
  syncExternalToForm(newExternal: TExternal, force = false): void {
    const formValues = this.options.getFormValues()
    const expectedFormValues = this.options.externalToForm(newExternal)

    // 检查是否需要同步
    if (!force && isEqual(formValues, expectedFormValues)) {
      return
    }

    // 检查是否是自己保存后的回显
    if (this.lastSavedFormValuesRef && isEqual(expectedFormValues, this.lastSavedFormValuesRef)) {
      console.log('📌 检测到回显数据，跳过同步')
      this.externalDataRef = newExternal
      return
    }

    console.log('📥 外部数据同步到表单')
    this.mutex.withLock(() => {
      this.options.setFormValues(expectedFormValues)
      this.externalDataRef = newExternal
    })
  }

  /**
   * 处理表单变更（在表单的 onChange 中调用）
   */
  handleFormChange(): void {
    if (this.mutex.isLocked) {
      console.log('🔒 系统更新，跳过保存')
      return
    }

    console.log('✏️ 用户操作，触发保存')
    const formValues = this.options.getFormValues()
    const externalValues = this.options.formToExternal(formValues, this.externalDataRef)

    this.lastSavedFormValuesRef = formValues
    this.options.onSave(externalValues)
  }

  /**
   * 重置控制器状态
   */
  reset(): void {
    this.mutex.reset()
    this.lastSavedFormValuesRef = null
  }
}
```

---

## 模式 4：使用泛型装饰器模式

为任何对象的方法添加互斥锁。

```typescript
/**
 * 方法装饰器 - 为类方法添加互斥锁
 * @example
 * class MyForm {
 *   private mutex = new MutexWrapper();
 *
 *   @WithMutex('mutex')
 *   setValues(values: any) {
 *     // 这个方法执行时会自动加锁
 *   }
 *
 *   onChange() {
 *     if (this.mutex.isLocked) return;
 *     this.save();
 *   }
 * }
 */
export function WithMutex(mutexPropertyName: string) {
  return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value

    descriptor.value = function (...args: any[]) {
      const mutex = this[mutexPropertyName] as MutexWrapper
      if (!mutex || !(mutex instanceof MutexWrapper)) {
        throw new Error(`${mutexPropertyName} is not a MutexWrapper instance`)
      }

      return mutex.withLock(() => originalMethod.apply(this, args))
    }

    return descriptor
  }
}
```

---

## 🎯 实战应用：重构你的代码

### 使用 FormSyncController 重构

```typescript
import { useMemo, useRef, useEffect } from 'react'
import { FormSyncController } from '@/utils/FormSyncController'

export const Transformer = (props: TransformerProps) => {
  const { codeDetail, onSave } = props

  // 创建表单实例（只在 ID 变化时重建）
  const formInstance = useMemo(() => {
    return createForm({
      initialValues: transferSchemaToForm(codeDetail)
    })
  }, [codeDetail.id])

  // 创建同步控制器
  const syncController = useRef<FormSyncController<LanderQuerySchemaDetail, any>>()

  if (!syncController.current) {
    syncController.current = new FormSyncController(codeDetail, {
      externalToForm: transferSchemaToForm,
      formToExternal: transferFormToSchema,
      setFormValues: values => formInstance.setValues(values),
      getFormValues: () => formInstance.getValuesIn(''),
      onSave: values => onSave(codeDetail.id, values)
    })
  }

  // 监听表单变化
  useEffect(() => {
    const dispose = formInstance.onFormValuesChange(() => {
      syncController.current?.handleFormChange()
    })
    return dispose
  }, [formInstance])

  // 同步外部数据
  useEffect(() => {
    syncController.current?.syncExternalToForm(codeDetail)
  }, [codeDetail])

  return (
    <FormProvider form={formInstance}>
      <Form>
        <TransformerComputeRender onSave={onSave} codeDetail={codeDetail} />
        <JSTransformer codeDetail={codeDetail} onSave={onSave} />
      </Form>
    </FormProvider>
  )
}
```

### 使用 MutexWrapper 重构（轻量级）

```typescript
import { MutexWrapper } from '@/utils/MutexWrapper'

export const Transformer = (props: TransformerProps) => {
  const { codeDetail, onSave } = props
  const mutexRef = useRef(new MutexWrapper())
  const detailRef = useRef(codeDetail)

  useEffect(() => {
    detailRef.current = codeDetail
  }, [codeDetail])

  const formInstance = useMemo(() => {
    return createForm({
      initialValues: transferSchemaToForm(codeDetail),
      effects: () => {
        onFormValuesChange(form => {
          // 🔑 简洁的互斥检查
          if (mutexRef.current.isLocked) {
            return
          }

          const formValues = form.getValuesIn('')
          const result = transferFormToSchema(formValues, detailRef.current)
          onSave?.(codeDetail.id, result)
        })
      }
    })
  }, [codeDetail.id])

  useEffect(() => {
    const formValues = formInstance.getValuesIn('')
    const expectedValues = transferSchemaToForm(codeDetail)

    if (!isEqual(formValues, expectedValues)) {
      // 🔑 使用 withLock 自动管理锁
      mutexRef.current.withLock(() => {
        formInstance.setValues(expectedValues)
      })
    }
  }, [codeDetail, formInstance])

  return <FormProvider form={formInstance}>{/* ...existing code... */}</FormProvider>
}
```

---

## 📦 完整工具包导出

```typescript
export { MutexWrapper } from './MutexWrapper'
export { FormSyncController } from './FormSyncController'
export { WithMutex } from './withMutex'
export type { FormSyncOptions } from './FormSyncController'

// React Hooks
export { useMutexState } from '../hooks/useMutexState'
export type { MutexStateReturn } from '../hooks/useMutexState'
```

---

## 🎓 总结：何时用哪种模式

| 模式                   | 适用场景               | 优点                 | 缺点                   |
| ---------------------- | ---------------------- | -------------------- | ---------------------- |
| **MutexWrapper**       | 简单的加锁需求         | 轻量、灵活           | 需要手动管理           |
| **useMutexState**      | React 状态需要区分来源 | 开箱即用             | 仅限 React             |
| **FormSyncController** | 复杂表单双向同步       | 功能完整、自动防循环 | 较重，需要配置         |
| **WithMutex 装饰器**   | 类方法需要自动加锁     | 声明式、优雅         | 仅限 Class 组件/工具类 |
