type FunctionComponent = (props: Record<PropertyKey, any>) => VirtualDom

import type { Children, IProps, VirtualDom } from './typings'

/**
 * @param { string | FunctionComponent } type - valid HTML tag name or Function Component
 * @param { Record<PropertyKey, any> } [props] - properties.
 * @param { ...Children} [children] - elements as rest arguments
 * @return { VirtualDom }
 */
function createElement(
  type: string | FunctionComponent,
  props: Record<PropertyKey, any>,
  ...children: Children[]
): VirtualDom {
  if (typeof type === 'function') return type({ children, ...props })

  return {
    type,
    props: {
      ...props,
      children,
    },
  }
}

/**
 * @param { VirtualDom }
 * @returns { HTMLElement }
 */
function render(json: VirtualDom | string): HTMLElement | Text {
  if (typeof json === 'string') return document.createTextNode(json)

  const {
    type,
    props: { children, ...attrs },
  } = json

  const node = document.createElement(type)
  for (const [attr, value] of Object.entries(attrs)) {
    // @ts-ignore
    node[attr] = value
  }

  const childrenArr = Array.isArray(children) ? children : [children]
  for (let child of childrenArr) {
    node.append(render(child))
  }

  return node
}

if (require.main === module) {
  const h = createElement
  const Container: FunctionComponent = ({ children, ...res }) => h('div', res, ...children)
  const Title: FunctionComponent = ({ children, ...res }) => h('b', res, ...children)
  const Paragraph: FunctionComponent = ({ children, ...res }) => h('p', res, ...children)
  const Button: FunctionComponent = ({ children, ...res }) => h('button', res, ...children)
  const Link: FunctionComponent = ({ children, ...res }) => h('a', res, ...children)
  const Name: FunctionComponent = ({ children, ...res }) => h('a', res, ...children)
  h(
    Container,
    {},
    h(Title, {}, ' this is '),
    h(
      Paragraph,
      { className: 'paragraph' },
      ' a ',
      h(Button, {}, ' button '),
      ' from ',
      h(Link, { href: 'https://bfe.dev' }, h(Name, {}, 'BFE'), '.dev')
    )
  )
}
