type JSXOpeningElement = {
  tag: string
}

type JSXClosingElement = {
  tag: string
}

type JSXChildren = string[]

type JSXElement = {
  openingElement: JSXOpeningElement
  children: JSXChildren
  closingElement: JSXClosingElement
}

function parse(code: string): JSXElement {
  // your code here
}

function generate(ast: JSXElement): string {
  // your code here
}
