## useControllableValue

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

```ts
export default function useCreation<T>(factory: () => T, deps: DependencyList): T
```

## useEventEmitter

```ts
type Subscription<T> = (val: T) => void
export declare class EventEmitter<T> {
  private subscriptions
  emit: (val: T) => void
  useSubscription: (callback: Subscription<T>) => void
}
export default function useEventEmitter<T = void>(): EventEmitter<T>
```

## useIsomorphicLayoutEffect

```ts
declare const useIsomorphicLayoutEffect: typeof useEffect
```

## useLatest

```ts
declare function useLatest<T>(value: T): import('react').MutableRefObject<T>
```

## useMemoizedFn

```ts
type noop = (this: any, ...args: any[]) => any
declare function useMemoizedFn<T extends noop>(fn: T): T
```

## useReactive

```ts
declare function useReactive<S extends Record<string, any>>(initialState: S): S
```
