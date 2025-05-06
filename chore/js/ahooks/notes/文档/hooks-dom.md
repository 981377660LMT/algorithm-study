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
declare function useEventListener(
  eventName: string | string[],
  handler: (event: Event) => void,
  options?: Options<Window>
): void
declare function useEventListener(
  eventName: string | string[],
  handler: noop,
  options: Options
): void
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

监听目标元素外的点击事件。

```ts
type DocumentEventKey = keyof DocumentEventMap
export default function useClickAway<T extends Event = Event>(
  onClickAway: (event: T) => void,
  target: BasicTarget | BasicTarget[],
  eventName?: DocumentEventKey | DocumentEventKey[]
): void
```

```tsx
import React, { useState, useRef } from 'react'
import { useClickAway } from 'ahooks'

export default () => {
  const [counter, setCounter] = useState(0)
  const ref = useRef(null)
  useClickAway(
    () => {
      setCounter(s => s + 1)
    },
    ref,
    ['click', 'contextmenu']
  )

  return (
    <div>
      <button type="button" ref={ref}>
        box
      </button>
      <p>counter: {counter}</p>
    </div>
  )
}
```

- 支持多个 DOM 对象
- 支持监听其它事件
  通过设置 eventName，可以指定需要监听的事件，试试点击鼠标右键。
- 支持 shadow DOM

## useDocumentVisibility

监听页面是否可见.

```ts
type VisibilityState = 'hidden' | 'visible' | 'prerender' | undefined
declare function useDocumentVisibility(): VisibilityState
```

```tsx
import React, { useEffect } from 'react'
import { useDocumentVisibility } from 'ahooks'

export default () => {
  const documentVisibility = useDocumentVisibility()
  useEffect(() => {
    console.log(`Current document visibility state: ${documentVisibility}`)
  }, [documentVisibility])
  return <div>Current document visibility state: {documentVisibility}</div>
}
```

## useDrop & useDrag

处理元素拖拽的 Hook。
useDrop 可以单独使用来接收文件、文字和网址的拖拽 => 文件上传场景。
useDrag 允许一个 DOM 节点被拖拽，需要配合 useDrop 使用。
向节点内触发粘贴动作也会被视为拖拽。

```ts
export interface DropOptions {
  onFiles?: (files: File[], event?: React.DragEvent) => void
  onUri?: (url: string, event?: React.DragEvent) => void
  onDom?: (content: any, event?: React.DragEvent) => void
  onText?: (text: string, event?: React.ClipboardEvent) => void

  onDragEnter?: (event?: React.DragEvent) => void
  onDragOver?: (event?: React.DragEvent) => void
  onDragLeave?: (event?: React.DragEvent) => void

  onDrop?: (event?: React.DragEvent) => void
  onPaste?: (event?: React.ClipboardEvent) => void
}
declare const useDrop: (target: BasicTarget, options?: DropOptions) => void

export interface DragOptions {
  onDragStart?: (event: React.DragEvent) => void
  onDragEnd?: (event: React.DragEvent) => void

  // 自定义拖拽过程中跟随鼠标指针的图像,<img> 或者 <canvas> 元素
  dragImage?: {
    image: string | Element
    offsetX?: number
    offsetY?: number
  }
}
declare const useDrag: <T>(
  data: T /** 拖拽的内容 **/,
  target: BasicTarget,
  options?: DragOptions
) => void
```

```tsx
import React, { useRef, useState } from 'react'
import { useDrop, useDrag } from 'ahooks'

const DragItem = ({ data }) => {
  const dragRef = useRef(null)
  const [dragging, setDragging] = useState(false)
  useDrag(data, dragRef, {
    onDragStart: () => {
      setDragging(true)
    },
    onDragEnd: () => {
      setDragging(false)
    }
  })

  return (
    <div
      ref={dragRef}
      style={{
        border: '1px solid #e8e8e8',
        padding: 16,
        width: 80,
        textAlign: 'center',
        marginRight: 16
      }}
    >
      {dragging ? 'dragging' : `box-${data}`}
    </div>
  )
}

export default () => {
  const [isHovering, setIsHovering] = useState(false)

  const dropRef = useRef(null)
  useDrop(dropRef, {
    onText: (text, e) => {
      console.log(e)
      alert(`'text: ${text}' dropped`)
    },
    onFiles: (files, e) => {
      console.log(e, files)
      alert(`${files.length} file dropped`)
    },
    onUri: (uri, e) => {
      console.log(e)
      alert(`uri: ${uri} dropped`)
    },
    onDom: (content: string, e) => {
      alert(`custom: ${content} dropped`)
    },
    onDragEnter: () => setIsHovering(true),
    onDragLeave: () => setIsHovering(false)
  })

  return (
    <div>
      <div ref={dropRef} style={{ border: '1px dashed #e8e8e8', padding: 16, textAlign: 'center' }}>
        {isHovering ? 'release here' : 'drop here'}
      </div>

      <div style={{ display: 'flex', marginTop: 8, overflow: 'auto' }}>
        {['1', '2', '3', '4', '5'].map(e => (
          <DragItem key={e} data={e} />
        ))}
      </div>
    </div>
  )
}
```

本质区别：

- **useDrag**：让元素**可以被拖动**，用于`拖拽源`。你用它让某个组件变成“可拖动的对象”。
- **useDrop**：让元素**可以接收拖动放下**，用于`放置目标`。你用它让某个组件变成“拖拽的接收区”。

**一句话总结**：  
useDrag 负责“拖”，useDrop 负责“接”。

## useEventTarget

常见表单控件(通过 e.target.value 获取表单值) 的 onChange 跟 value 逻辑封装，**支持自定义值转换和重置功能**。
onChange、reset。

```ts
interface EventTarget<U> {
  target: { value: U }
}
export interface Options<T, U> {
  initialValue?: T
  transformer?: (value: U) => T // U 是原始组件的值，T 是业务需要的值
}
declare function useEventTarget<T, U = T>(
  options?: Options<T, U>
): readonly [
  T | undefined,
  {
    readonly onChange: (e: EventTarget<U>) => void
    readonly reset: () => void
  }
]
```

```tsx
import React from 'react'
import { useEventTarget } from 'ahooks'

export default () => {
  const [value, { onChange, reset }] = useEventTarget({
    initialValue: '',
    transformer: (val: string) => val.replace(/[^\d]/g, '')
  })

  return (
    <div>
      <input
        value={value}
        onChange={onChange}
        style={{ width: 200, marginRight: 20 }}
        placeholder="Please type here"
      />
      <button type="button" onClick={reset}>
        reset
      </button>
    </div>
  )
}
```

## useExternal

动态注入 JS 或 CSS 资源，useExternal 可以`保证资源全局唯一`。

```ts
type JsOptions = {
  type: 'js'
  js?: Partial<HTMLScriptElement>
  keepWhenUnused?: boolean
}
type CssOptions = {
  type: 'css'
  css?: Partial<HTMLStyleElement>
  keepWhenUnused?: boolean
}
type DefaultOptions = {
  type?: never // 支持 js/css，如果不传，则根据 path 推导
  js?: Partial<HTMLScriptElement> // script 标签支持的属性
  css?: Partial<HTMLStyleElement> // link 标签支持的属性
  keepWhenUnused?: boolean // 在不持有资源的引用后，仍然保留资源，默认 false
}
export type Options = JsOptions | CssOptions | DefaultOptions
export type Status = 'unset' | 'loading' | 'ready' | 'error' // 加载状态，unset(未设置), loading(加载中), ready(加载完成), error(加载失败)
declare const useExternal: (path?: string /** 外部资源地址 **/, options?: Options) => Status
```

```tsx
import { useExternal } from 'ahooks'
import React, { useState } from 'react'

export default () => {
  const [path, setPath] = useState('/useExternal/bootstrap-badge.css')

  const status = useExternal(path)

  return (
    <>
      <p>
        Status: <b>{status}</b>
      </p>
      <div className="bd-example" style={{ wordBreak: 'break-word' }}>
        <span className="badge badge-pill badge-primary">Primary</span>
        <span className="badge badge-pill badge-secondary">Secondary</span>
        <span className="badge badge-pill badge-success">Success</span>
        <span className="badge badge-pill badge-danger">Danger</span>
        <span className="badge badge-pill badge-warning">Warning</span>
        <span className="badge badge-pill badge-info">Info</span>
        <span className="badge badge-pill badge-light">Light</span>
        <span className="badge badge-pill badge-dark">Dark</span>
      </div>
      <br />
      <button type="button" style={{ marginRight: 8 }} onClick={() => setPath('')}>
        unload
      </button>
      <button
        type="button"
        style={{ marginRight: 8 }}
        onClick={() => setPath('/useExternal/bootstrap-badge.css')}
      >
        load
      </button>
    </>
  )
}
```

## useTitle

```ts
export interface Options {
  restoreOnUnmount?: boolean // 组件卸载时，是否恢复上一个页面标题，默认false
}
declare function useTitle(title: string, options?: Options): void
```

## useFavicon

```ts
// href：favicon 地址, 支持 svg/png/ico/gif 后缀的图片
declare const useFavicon: (href: string) => void
```

## useFullscreen

管理 DOM 全屏的 Hook。

- dom 元素全屏
- 图片全屏
- 页面全屏
- 与其它全屏操作共存

```ts
export interface PageFullscreenOptions {
  className?: string
  zIndex?: number
}
export interface Options {
  onEnter?: () => void
  onExit?: () => void
  pageFullscreen?: boolean | PageFullscreenOptions // 是否是页面全屏。当参数类型为对象时，可以设置全屏元素的类名和 z-index
}
declare const useFullscreen: (
  target: BasicTarget,
  options?: Options
) => readonly [
  boolean /** isFullScreen **/,
  {
    readonly enterFullscreen: () => void
    readonly exitFullscreen: () => void
    readonly toggleFullscreen: () => void
    readonly isEnabled: true // 是否支持全屏
  }
]
```

## useHover

```ts
export interface Options {
  onEnter?: () => void
  onLeave?: () => void
  onChange?: (isHovering: boolean) => void
}
declare const _default: (target: BasicTarget, options?: Options) => boolean
```

## useMutationObserver

一个监听指定的 DOM 树发生变化的 Hook

```ts
declare const useMutationObserver: (
  callback: MutationCallback,
  target: BasicTarget,
  options?: MutationObserverInit
) => void
```

```tsx
import { useMutationObserver } from 'ahooks'
import React, { useRef, useState } from 'react'

const App: React.FC = () => {
  const [width, setWidth] = useState(200)
  const [count, setCount] = useState(0)

  const ref = useRef<HTMLDivElement>(null)

  useMutationObserver(
    mutationsList => {
      mutationsList.forEach(() => setCount(c => c + 1))
    },
    ref,
    { attributes: true }
  )

  return (
    <div>
      <div ref={ref} style={{ width, padding: 12, border: '1px solid #000', marginBottom: 8 }}>
        current width：{width}
      </div>
      <button onClick={() => setWidth(w => w + 10)}>widening</button>
      <p>Mutation count {count}</p>
    </div>
  )
}

export default App
```

## useInViewport

观察元素是否在可见区域，以及元素可见比例。

```ts
type CallbackType = (entry: IntersectionObserverEntry) => void
export interface Options {
  rootMargin?: string
  threshold?: number | number[]
  root?: BasicTarget<Element>
  callback?: CallbackType
}
declare function useInViewport(
  target: BasicTarget | BasicTarget[],
  options?: Options
): readonly [boolean | undefined, number | undefined]
```

## useKeyPress

```ts
export type KeyType = number | string
export type KeyPredicate = (event: KeyboardEvent) => KeyType | boolean | undefined
export type KeyFilter = KeyType | KeyType[] | ((event: KeyboardEvent) => boolean)
export type KeyEvent = 'keydown' | 'keyup'
export type Target = BasicTarget<HTMLElement | Document | Window>
export type Options = {
  events?: KeyEvent[]
  target?: Target
  exactMatch?: boolean
  useCapture?: boolean
}
declare function useKeyPress(
  keyFilter: KeyFilter,
  eventHandler: (event: KeyboardEvent, key: KeyType) => void,
  option?: Options
): void
```

## useLongPress

```ts
type EventType = MouseEvent | TouchEvent
export interface Options {
  delay?: number
  moveThreshold?: {
    x?: number
    y?: number
  }
  onClick?: (event: EventType) => void
  onLongPressEnd?: (event: EventType) => void
}
declare function useLongPress(
  onLongPress: (event: EventType) => void,
  target: BasicTarget,
  { delay, moveThreshold, onClick, onLongPressEnd }?: Options
): void
```

## useMouse

```ts
export interface CursorState {
  screenX: number
  screenY: number
  clientX: number
  clientY: number
  pageX: number
  pageY: number
  elementX: number
  elementY: number
  elementH: number
  elementW: number
  elementPosX: number
  elementPosY: number
}
declare const _default: (target?: BasicTarget) => CursorState
```

## useResponsive

```ts
type ResponsiveConfig = Record<string, number>
type ResponsiveInfo = Record<string, boolean>
export declare function configResponsive(config: ResponsiveConfig): void
declare function useResponsive(): ResponsiveInfo
```

## useScroll

```ts
type Position = {
  left: number
  top: number
}
export type Target = BasicTarget<Element | Document>
export type ScrollListenController = (val: Position) => boolean
declare function useScroll(
  target?: Target,
  shouldUpdate?: ScrollListenController
): Position | undefined
```

## useSize

```ts
type Size = {
  width: number
  height: number
}
declare function useSize(target: BasicTarget): Size | undefined
```

## useFocusWithin

```ts
export interface Options {
  onFocus?: (e: FocusEvent) => void
  onBlur?: (e: FocusEvent) => void
  onChange?: (isFocusWithin: boolean) => void
}
export default function useFocusWithin(target: BasicTarget, options?: Options): boolean
```

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
