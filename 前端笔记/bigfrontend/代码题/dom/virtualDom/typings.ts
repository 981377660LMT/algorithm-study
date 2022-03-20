interface VirtualDom {
  type: keyof HTMLElementTagNameMap | (string & {})
  props: IProps
}

interface IProps {
  children: Children[]
  [attr: string]: any
}

type Children = VirtualDom | string

export { VirtualDom, Children, IProps }
