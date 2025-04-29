## useEventListener

优雅的使用 addEventListener。
可以看到重载、类型推导都很完美。

```ts
type noop = (...p: any) => void
export type Target = BasicTarget<HTMLElement | Element | Window | Document>
type Options<T extends Target = Target> = {
  target?: T
  capture?: boolean
  once?: boolean
  passive?: boolean // 设置为 true 时，表示 listener 永远不会调用 preventDefault()
  enable?: boolean // 是否开启监听
}
declare function useEventListener<K extends keyof HTMLElementEventMap>(
  eventName: K,
  handler: (ev: HTMLElementEventMap[K]) => void,
  options?: Options<HTMLElement>
): void
declare function useEventListener<K extends keyof ElementEventMap>(
  eventName: K,
  handler: (ev: ElementEventMap[K]) => void,
  options?: Options<Element>
): void
declare function useEventListener<K extends keyof DocumentEventMap>(
  eventName: K,
  handler: (ev: DocumentEventMap[K]) => void,
  options?: Options<Document>
): void
declare function useEventListener<K extends keyof WindowEventMap>(
  eventName: K,
  handler: (ev: WindowEventMap[K]) => void,
  options?: Options<Window>
): void
declare function useEventListener(eventName: string | string[], handler: (event: Event) => void, options?: Options<Window>): void
declare function useEventListener(eventName: string | string[], handler: noop, options: Options): void
```

```tsx
import React, { useRef, useState } from 'react'
import { useEventListener } from 'ahooks'

export default () => {
  const ref = useRef(null)
  const [value, setValue] = useState('')

  useEventListener(
    ['mouseenter', 'mouseleave'],
    ev => {
      setValue(ev.type)
    },
    { target: ref }
  )

  return (
    <button ref={ref} type="button">
      You Option is {value}
    </button>
  )
}
```

## useClickAway

## useDocumentVisibility

## useDrop & useDrag

## useEventTarget

## useExternal

## useTitle

## useFavicon

## useFullscreen

## useHover

## useMutationObserver

## useInViewport

## useKeyPress

## useLongPress

## useMouse

## useResponsive

## useScroll

## useSize

## useFocusWithin

---

`passive:true`

- 是什么
  passive: true 表示事件处理函数不会调用 event.preventDefault().
  passive: true 是给事件监听器（如 addEventListener）用的一个参数，主要用于提升滚动等高频事件的性能.
  用于滚动相关事件（如 touchstart、touchmove、wheel、scroll）.

  ```ts
  // 推荐：监听滚动时加 passive:true
  window.addEventListener('scroll', handler, { passive: true })

  // 不推荐：如果你要阻止默认行为
  window.addEventListener('touchmove', e => e.preventDefault(), { passive: false })
  ```

  - 为什么要叫做 passive?
    之所以叫做 passive，**是因为事件监听器是“被动的”，不会主动干预（阻止）浏览器的默认行为（如滚动）。**
    也就是说，监听器只是“被动地”接收事件，不会通过 preventDefault() 阻止事件的默认动作。
    这样浏览器可以放心优化，提高性能。

- 为什么
  **如果你不加 passive: true，浏览器每次触发事件都要等你的回调执行完，看看你会不会 preventDefault()，才能决定要不要滚动页面。**
  **加了 passive: true，浏览器不用等，直接滚动，滚动更流畅，不卡顿。**

- 怎么办
  passive: true 让浏览器知道你的事件处理不会阻止默认行为，从而优化滚动等高频事件的性能。

---

```ts
type Dispatch<T> = (value: T) => void
type SetAction<S> = S | ((prevState: S) => S)
```
