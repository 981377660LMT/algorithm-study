// https://pomb.us/build-your-own-react/

// import { createElement } from 'react'

// const a = createElement('a')
//  返回DetailedReactHTMLElement
//  常用的：
// const element = { type: 'h1', props: { title: 'foo', children: 'Hello' } }

function createElement(type: any, props: any, ...children: any[]) {
  return {
    type,
    props: {
      ...props,
      children: children.map(child =>
        typeof child === 'object' ? child : createTextElement(child)
      ),
    },
  }
}

function createTextElement(text: any) {
  return {
    type: 'TEXT_ELEMENT',
    props: {
      nodeValue: text,
      children: [],
    },
  }
}

function render(
  element: { type: string; props: { [x: string]: any; children?: any } },
  container: { appendChild: (arg0: any) => void }
) {
  const dom =
    element.type == 'TEXT_ELEMENT'
      ? document.createTextNode('')
      : document.createElement(element.type)
  const isProperty = (key: string) => key !== 'children'
  Object.keys(element.props)
    .filter(isProperty)
    .forEach(name => {
      Reflect.set(dom, name, element.props[name])
    })
  element.props.children.forEach((child: any) => render(child, dom))
  container.appendChild(dom)
}

const Didact = {
  createElement,
  render,
}
/** @jsx Didact.createElement */
const element = (
  <div id="foo">
    <a>bar</a>
    <b />
  </div>
)
const container = document.getElementById('root')!
Didact.render(element, container)
