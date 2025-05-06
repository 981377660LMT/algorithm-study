## useControllableValue

在某些组件开发时，我们需要组件的状态既可以自己管理，也可以被外部控制，useControllableValue 就是帮你管理这种状态的 Hook。

- 如果 props 中没有 value，则组件内部自己管理 state
- 如果 props 有 value 字段，则由父级接管控制 state

```ts
export interface Options<T> {
  defaultValue?: T
  defaultValuePropName?: string
  valuePropName?: string
  trigger?: string
}
export type Props = Record<string, any>
export interface StandardProps<T> {
  value: T
  defaultValue?: T
  onChange: (val: T) => void
}
declare function useControllableValue<T = any>(
  props: StandardProps<T>
): [T, (v: SetStateAction<T>) => void]
declare function useControllableValue<T = any>(
  props?: Props,
  options?: Options<T>
): [T, (v: SetStateAction<T>, ...args: any[]) => void]
```

## useCreation

useCreation 是 useMemo 或 useRef 的替代品。
因为 useMemo 不能保证被 memo 的值一定不会被重新计算，而 useCreation 可以保证这一点。

```ts
export default function useCreation<T>(factory: () => T, deps: DependencyList): T
```

```ts
const a = useRef(new Subject()) // 每次重渲染，都会执行实例化 Subject 的过程，即便这个实例立刻就被扔掉了
const b = useCreation(() => new Subject(), []) // 通过 factory 函数，可以避免性能隐患
```

## useEventEmitter

在组件多次渲染时，每次渲染调用 useEventEmitter 得到的返回值会保持不变，不会重复创建 EventEmitter 的实例。

```ts
type Subscription<T> = (val: T) => void
export declare class EventEmitter<T> {
  private subscriptions
  emit: (val: T) => void
  useSubscription: (callback: Subscription<T>) => void // useSubscription 会在组件创建时自动注册订阅，并在组件销毁时自动取消订阅。
}
export default function useEventEmitter<T = void>(): EventEmitter<T>
```

对于子组件通知父组件的情况，我们仍然推荐直接使用 `props 传递一个 onEvent 函数`。
而对于父组件通知子组件的情况，可以使用 `forwardRef` 获取子组件的 ref ，再进行子组件的方法调用。 useEventEmitter 适合的是在距离较远的组件之间进行事件通知，或是在多个组件之间共享事件通知。

## useIsomorphicLayoutEffect

在非浏览器环境返回 useEffect，在浏览器环境返回 useLayoutEffect

```ts
declare const useIsomorphicLayoutEffect: typeof useEffect
```

## useLatest

```ts
declare function useLatest<T>(value: T): import('react').MutableRefObject<T>
```

## useMemoizedFn

在某些场景中，我们需要使用 useCallback 来记住一个函数，但是在第二个参数 deps 变化时，会重新生成函数，导致函数地址变化。
使用 useMemoizedFn，可以省略第二个参数 deps，同时保证函数地址永远不会变化。

- 解决闭包问题
- 性能优化

```ts
type noop = (this: any, ...args: any[]) => any
declare function useMemoizedFn<T extends noop>(fn: T): T
```

注意：useMemoizedFn 持久后的函数不会继承函数本身的属性

## useReactive

提供一种数据响应式的操作体验，定义数据状态不需要写useState，直接修改属性即可刷新视图。
useReactive 产生可操作的代理对象一直都是同一个引用，useEffect , useMemo ,useCallback ,子组件属性传递 等如果依赖的是这个代理对象是不会引起重新执行。

```ts
declare function useReactive<S extends Record<string, any>>(initialState: S): S
```
