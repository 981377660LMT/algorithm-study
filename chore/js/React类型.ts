// https://github.com/981377660LMT/ts/issues/842

// ReactNode => 最广泛，适用于 {xxx}，表示任意可渲染的元素；
// ReactElement => JSX，表示React.createElement 创建的元素；
// ComponentType => 严格要求类组件 or 函数组件，适用于 <Comp />
// ElementType => 兼容 ComponentType 和 普通HTML/SVG 元素(JSX.IntrinsicElements)，适用于 <Comp />

type ElementType<P = any> =
  | {
      [K in keyof JSX.IntrinsicElements]: P extends JSX.IntrinsicElements[K] ? K : never
    }[keyof JSX.IntrinsicElements]
  | ComponentType<P>

type ComponentType<P = {}> = ComponentClass<P> | FunctionComponent<P>

// 构造签名或调用签名
// !历史原因导致两种返回不一样，且调用签名手动声明 null
// https://stackoverflow.com/questions/58123398/when-to-use-jsx-element-vs-reactnode-vs-reactelement
type JSXElementConstructor<P> =
  | ((props: P) => ReactElement<any, any> | null)
  | (new (props: P) => Component<any, any>)

interface ReactElement<
  P = any,
  T extends string | JSXElementConstructor<any> = string | JSXElementConstructor<any>
> {
  type: T
  props: P
  key: string | null
}

interface ComponentClass<P = {}, S = ComponentState> extends StaticLifecycle<P, S> {
  new (props: P, context?: any): Component<P, S>
  propTypes?: WeakValidationMap<P> | undefined
  contextType?: Context<any> | undefined
  contextTypes?: ValidationMap<any> | undefined
  childContextTypes?: ValidationMap<any> | undefined
  defaultProps?: Partial<P> | undefined
  displayName?: string | undefined
}

interface FunctionComponent<P = {}> {
  (props: PropsWithChildren<P>, context?: any): ReactElement<any, any> | null
  propTypes?: WeakValidationMap<P> | undefined
  contextTypes?: ValidationMap<any> | undefined
  defaultProps?: Partial<P> | undefined
  displayName?: string | undefined
}

type ReactNode = ReactChild | ReactFragment | ReactPortal | boolean | null | undefined

namespace JSX {
  interface Element extends React.ReactElement<any, any> {}
}

export {}
