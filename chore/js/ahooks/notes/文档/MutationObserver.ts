/* eslint-disable no-var */
/* eslint-disable vars-on-top */
/**
 * MutationObserver API 的回调函数类型。
 * 当观察到 DOM 树发生变化时调用。
 *
 * @param mutations 一个 MutationRecord 对象数组，描述了每个发生的 DOM 变化。
 * @param observer 调用此回调的 MutationObserver 实例。
 */
type MutationCallback = (mutations: MutationRecord[], observer: MutationObserver) => void

/**
 * MutationObserverInit 对象用于配置 MutationObserver 实例的行为。
 */
interface MutationObserverInit {
  /**
   * 设置为 true 以观察目标节点属性的变化。
   * 默认为 false。
   */
  attributes?: boolean
  /**
   * 如果 attributes 设置为 true，并且希望在回调中记录属性变化前的值，则设置为 true。
   * 默认为 false。
   */
  attributeOldValue?: boolean
  /**
   * 如果 attributes 设置为 true，并且希望只观察特定属性的变化，则设置为一个包含属性名称（不含命名空间）的数组。
   * 如果未指定，则观察所有属性的变化。
   */
  attributeFilter?: string[]

  /**
   * 设置为 true 以观察目标节点字符数据的变化（例如文本节点的内容）。
   * 默认为 false。
   */
  characterData?: boolean
  /**
   * 如果 characterData 设置为 true，并且希望在回调中记录字符数据变化前的值，则设置为 true。
   * 默认为 false。
   */
  characterDataOldValue?: boolean

  /**
   * 设置为 true 以观察目标节点的子节点（直接子元素）的添加或删除。
   * 默认为 false。
   */
  childList?: boolean
  /**
   * 设置为 true 以将观察扩展到目标节点的整个子树。
   * 如果设置为 true，则对目标节点子树中的所有节点（包括目标节点自身）的更改都会被观察。
   * childList、attributes 或 characterData 中至少有一个必须为 true 才能使 subtree 生效。
   * 默认为 false。
   */
  subtree?: boolean
}

/**
 * MutationRecord 描述了一个单独的 DOM 变化。
 * MutationObserver 的回调函数会接收一个 MutationRecord 对象的数组。
 */
interface MutationRecord {
  /**
   * 对于 "attributes" 类型的变化，是被修改属性的局部名称，否则为 null。
   */
  readonly attributeName: string | null
  /**
   * 对于 "attributes" 类型的变化，是被修改属性的命名空间，否则为 null。
   */
  readonly attributeNamespace: string | null
  /**
   * 返回一个 DOMString 数组，其中包含已添加的节点；如果未添加任何节点，则返回一个空的 NodeList。
   */
  readonly addedNodes: NodeList
  /**
   * 返回一个 DOMString 数组，其中包含已移除的节点；如果未移除任何节点，则返回一个空的 NodeList。
   */
  readonly removedNodes: NodeList
  /**
   * 返回变化前的属性值（如果 attributeOldValue 设置为 true），或变化前的字符数据（如果 characterDataOldValue 设置为 true），否则为 null。
   */
  readonly oldValue: string | null
  /**
   * 返回发生变化的节点之前的同级节点，如果不存在则为 null。
   */
  readonly previousSibling: Node | null
  /**
   * 返回发生变化的节点之后的同级节点，如果不存在则为 null。
   */
  readonly nextSibling: Node | null
  /**
   * 返回发生变化的节点。
   * 对于 "attributes" 类型的变化，这是属性被修改的元素。
   * 对于 "characterData" 类型的变化，这是字符数据被修改的 CharacterData 节点。
   * 对于 "childList" 类型的变化，这是子节点被添加或移除的节点。
   */
  readonly target: Node
  /**
   * 返回变化的类型。
   * "attributes": 属性值发生变化。
   * "characterData": CharacterData 节点的数据发生变化。
   * "childList": 子节点的添加或删除。
   */
  readonly type: 'attributes' | 'characterData' | 'childList'
}

/**
 * MutationObserver 接口提供了一种监视对 DOM 树所做更改的能力。
 * 它被设计为替换旧的 Mutation Events 功能。
 */
interface MutationObserver {
  /**
   * 停止 MutationObserver 实例接收通知，直到再次调用 observe() 方法。
   */
  disconnect(): void
  /**
   * 配置 MutationObserver 实例以开始接收有关给定目标节点上 DOM 变化的通知。
   * @param target 要观察 DOM 变化的 Node (例如 Element)。
   * @param options 一个 MutationObserverInit 对象，指定要报告哪些 DOM 变化。
   */
  observe(target: Node, options?: MutationObserverInit): void
  /**
   * 从 MutationObserver 的通知队列中移除所有待处理的记录，并将它们作为新的 MutationRecord 对象数组返回。
   * 这个方法最常见的用例是在断开观察器连接之前立即获取所有挂起的更改记录，以便可以处理任何未处理的更改。
   */
  takeRecords(): MutationRecord[]
}

/**
 * MutationObserver 构造函数。
 * @param callback 当观察到 DOM 变化时调用的函数。
 */
declare var MutationObserver: {
  prototype: MutationObserver
  new (callback: MutationCallback): MutationObserver
}

export {}
