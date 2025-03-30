// 1618. 找出适应屏幕的最大字号
// https://leetcode.cn/problems/maximum-font-to-fit-a-sentence-in-a-screen/description/

// 假设有 FontInfo 接口
interface FontInfo {
  getWidth(fontSize: number, ch: string): number
  getHeight(fontSize: number): number
}

// 示例中我们用 fontInfo 实例调用接口（实际面试中，可能由面试官提供）
declare const fontInfo: FontInfo

function findMaxFont(text: string, w: number, h: number, fonts: number[]): number {
  let left = 0
  let right = fonts.length - 1
  let res = -1

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

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    const fontSize = fonts[mid]
    if (check(fontSize)) {
      res = fontSize
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  return res
}

// 示例调用（这里给出伪代码，实际面试中 fontInfo 的实现由面试官提供）：
/*
const text = "helloworld";
const w = 80;
const h = 20;
const fonts = [6,8,10,12,14,16,18,24,36];
const maxFont = findMaxFont(text, w, h, fonts);
console.log(maxFont);  // 输出：6
*/

export {}
