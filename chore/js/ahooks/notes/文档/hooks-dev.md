## useTrackedEffect

```ts
type Effect<T extends DependencyList> = (
  changes?: number[],
  previousDeps?: T,
  currentDeps?: T
) => void | (() => void)
declare const useTrackedEffect: <T extends DependencyList>(effect: Effect<T>, deps?: [...T]) => void
```

## useWhyDidYouUpdate

```ts
export type IProps = Record<string, any>
export default function useWhyDidYouUpdate(componentName: string, props: IProps): void
```
