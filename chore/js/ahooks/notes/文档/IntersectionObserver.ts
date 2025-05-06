/* eslint-disable no-var */
/* eslint-disable vars-on-top */

/**
 * IntersectionObserver API 的回调函数类型。
 * 当一个或多个目标元素的交叉状态发生变化时调用。
 *
 * @param entries 一个 IntersectionObserverEntry 对象数组，描述了每个目标元素的交叉状态。
 * @param observer 调用此回调的 IntersectionObserver 实例。
 */
type IntersectionObserverCallback = (
  entries: IntersectionObserverEntry[],
  observer: IntersectionObserver
) => void

/**
 * IntersectionObserverInit 对象用于配置 IntersectionObserver 实例。
 */
interface IntersectionObserverInit {
  /**
   * 观察器的根元素。目标元素将与此元素的边界框进行比较。
   * 如果未指定或为 null，则默认为浏览器视口。
   */
  root?: Element | Document | null
  /**
   * 根元素的边距。可以用来扩大或缩小根元素的边界框，以便在计算交叉时进行调整。
   * 值的格式类似于 CSS margin 属性（例如 "10px 20px 30px 40px"）。百分比值无效。
   */
  rootMargin?: string
  /**
   * 一个数字或数字数组，表示触发回调的 intersectionRatio 阈值。
   * 当目标元素的可见比例达到或超过（或低于，取决于滚动方向）这些阈值时，回调将被调用。
   * 如果指定为数组，则每当 intersectionRatio 穿过数组中的任何值时，都会调用回调。
   * 默认值为 0（表示只要目标元素有任何像素可见，就会触发）。
   * 值为 1.0 表示只有当目标元素完全可见时才会触发。
   */
  threshold?: number | number[]
}

/**
 * IntersectionObserverEntry 描述了目标元素与其根容器在特定时间点的交叉状态。
 * IntersectionObserver 的回调函数会接收一个 IntersectionObserverEntry 对象的数组。
 */
interface IntersectionObserverEntry {
  /**
   * 目标元素的边界矩形，由 getBoundingClientRect() 返回。
   */
  readonly boundingClientRect: DOMRectReadOnly
  /**
   * 目标元素与根元素交叉区域的可见比例。
   * 范围从 0.0（完全不可见）到 1.0（完全可见）。
   */
  readonly intersectionRatio: number
  /**
   * 描述目标元素与根元素交叉区域的矩形。
   */
  readonly intersectionRect: DOMRectReadOnly
  /**
   * 一个布尔值，如果目标元素当前与根元素或文档视口交叉，则为 true。
   */
  readonly isIntersecting: boolean
  /**
   * 根元素的边界矩形。如果根是文档视口，则此值为 null。
   */
  readonly rootBounds: DOMRectReadOnly | null
  /**
   * 被观察的目标元素。
   */
  readonly target: Element
  /**
   * 从 IntersectionObserver 创建到交叉状态发生变化的时间戳（以毫秒为单位）。
   */
  readonly time: DOMHighResTimeStamp
}

/**
 * IntersectionObserver 接口提供了一种异步观察目标元素与其祖先元素或顶级文档视口（“根”）交叉状态变化的方法。
 */
interface IntersectionObserver {
  /**
   * 观察器的根元素。如果构造时未指定，则为 null（表示浏览器视口）。
   */
  readonly root: Element | Document | null
  /**
   * 一个字符串，其格式类似于 CSS margin 属性，用于在计算交叉之前扩展或收缩根元素的边界框。
   */
  readonly rootMargin: string
  /**
   * 一个数字或数字数组，表示触发回调的 intersectionRatio 阈值。
   */
  readonly thresholds: ReadonlyArray<number>
  /**
   * 停止观察所有目标元素。
   */
  disconnect(): void
  /**
   * 开始观察指定的目标元素。
   * @param target 要观察的目标元素。
   */
  observe(target: Element): void
  /**
   * 返回一个 IntersectionObserverEntry 对象数组，每个对象对应一个当前正在被观察的目标元素，
   * 无论它们是否与根元素交叉。
   */
  takeRecords(): IntersectionObserverEntry[]
  /**
   * 停止观察指定的目标元素。
   * @param target 要停止观察的目标元素。
   */
  unobserve(target: Element): void
}

/**
 * IntersectionObserver 构造函数。
 * @param callback 当目标元素的交叉状态发生变化时调用的函数。
 * @param options 一个可选的配置对象，用于自定义观察器的行为。
 */
declare var IntersectionObserver: {
  prototype: IntersectionObserver
  new (
    callback: IntersectionObserverCallback,
    options?: IntersectionObserverInit
  ): IntersectionObserver
}

export {}
