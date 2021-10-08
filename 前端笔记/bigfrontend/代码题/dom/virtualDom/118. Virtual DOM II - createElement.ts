/**
 * MyElement is the type your implementation supports
 *
 * type MyNode = MyElement | string
 */

/**
 * @param { string } type - valid HTML tag name
 * @param { object } [props] - properties.
 * @param { ...MyNode} [children] - elements as rest arguments
 * @return { MyElement }
 */
function createElement(type: string, props: object, ...children: MyNode[]): MyElement {
  // your code here
}

/**
 * @param { MyElement }
 * @returns { HTMLElement }
 */
function render(myElement): HTMLElement {
  // your code here
}
