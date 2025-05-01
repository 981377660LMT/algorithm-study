// 1618. 找出适应屏幕的最大字号
// https://leetcode.cn/problems/maximum-font-to-fit-a-sentence-in-a-screen/description/
//
// 给定一个字符串 text。并能够在 宽为 w 高为 h 的屏幕上显示该文本。
// 字体数组中包含按升序排列的可用字号，您可以从该数组中选择任何字体大小。
// 您可以使用FontInfo接口来获取任何可用字体大小的任何字符的宽度和高度。
// !请注意：文本最多只能排放一排
// 返回可用于在屏幕上显示文本的最大字体大小。如果文本不能以任何字体大小显示，则返回-1。

interface IFontInfo {
  getWidth(fontSize: number, ch: string): number
  getHeight(fontSize: number): number
}

declare const fontInfo: IFontInfo

function maxFont(text: string, w: number, h: number, fonts: number[]): number {
  if (!check(fonts[0])) return -1
  let left = 0
  let right = fonts.length - 1
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    const fontSize = fonts[mid]
    if (check(fontSize)) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return fonts[right]

  // 判断 fontSize 是否能在屏幕上显示 text
  function check(fontSize: number): boolean {
    if (fontInfo.getHeight(fontSize) > h) return false

    let totalWidth = 0
    for (let i = 0; i < text.length; i++) {
      totalWidth += fontInfo.getWidth(fontSize, text[i])
      if (totalWidth > w) return false
    }
    return true
  }
}

// 示例调用（这里给出伪代码，实际面试中 fontInfo 的实现由面试官提供）：
/*
const text = "helloworld";
const w = 80;
const h = 20;
const fonts = [6,8,10,12,14,16,18,24,36];
const maxFont = maxFont(text, w, h, fonts);
console.log(maxFont);  // 输出：6
*/

export {}
