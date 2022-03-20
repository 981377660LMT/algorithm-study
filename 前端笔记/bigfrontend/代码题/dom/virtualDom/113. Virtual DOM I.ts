interface VirtualDom {
  type: keyof HTMLElementTagNameMap | (string & {})
  props: {
    children?: Children[] | Children
    [attr: string]: any
  }
}

type Children = VirtualDom | string

/**
 * @param {Node}
 * @return {VirtualDom} object literal presentation
 */
function virtualize(element: HTMLElement | Text): VirtualDom | string {
  if (element instanceof Text) {
    return element.data
  } else {
    const res: VirtualDom = { type: element.tagName.toLowerCase(), props: {} }

    // props里的非children属性
    for (const attr of Array.from(element.attributes)) {
      const name = attr.name === 'class' ? 'className' : attr.name
      res.props[name] = attr.value
    }

    // props里children属性
    const children: (string | VirtualDom)[] = []
    for (const node of Array.from(element.childNodes)) {
      // text node  (1:元素节点，2：属性节点，废弃，3：文本节点)
      if (node.nodeType === 3) children.push(node.textContent!)
      else children.push(virtualize(node as HTMLElement))
    }

    res.props.children = children.length === 1 ? children[0] : children

    return res
  }
}

/**
 * @param {VirtualDom | string} valid object literal presentation
 * @return {Node}
 */
function render(json: VirtualDom | string): HTMLElement | Text {
  // 1.文本节点
  if (typeof json === 'string') return document.createTextNode(json)

  const {
    type,
    props: { children, ...attrs },
  } = json

  // 2.根节点的非children属性
  const element = document.createElement(type)
  for (const [attr, value] of Object.entries(attrs)) {
    element.setAttribute(attr, value)
  }

  // 3. 根节点的children
  const childrenArray = Array.isArray(children) ? children : [children]
  for (let child of childrenArray) {
    if (child == undefined) continue
    element.append(render(child))
  }

  return element
}

const json: VirtualDom = {
  type: 'div',
  props: {
    children: [
      {
        type: 'h1',
        props: {
          children: ' this is ',
        },
      },
      {
        type: 'p',
        props: {
          className: 'paragraph',
          children: [
            ' a ',
            {
              type: 'button',
              props: {
                children: ' button ',
              },
            },
            ' from',
            {
              type: 'a',
              props: {
                href: 'https://bfe.dev',
                children: [
                  {
                    type: 'b',
                    props: {
                      children: 'BFE',
                    },
                  },
                  '.dev',
                ],
              },
            },
          ],
        },
      },
    ],
  },
}

console.log(render(json))
console.log(virtualize(render(json)))

export {}
// @ts-ignore
// console.log(Array.from(undefined))
