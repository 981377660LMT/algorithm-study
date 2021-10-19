// 将一组表单元素编码为查询字符串。
// Encodes a set of form elements as a query string.
serializeForm(document.querySelector('#form')!)
// email=test%40email.com&name=Test%20Name

/**
 *
 * @param form
 * @returns
 * 1.表单转换为 FormData
 * 2.Array.prototype.map ()和 encodeURIComponent ()对每个字段的值进行编码。
 */
function serializeForm(form: HTMLFormElement): string {
  return Array.from<[key: string, value: string], string>(new FormData(form), entry =>
    entry.map(encodeURIComponent).join('=')
  ).join('&')
}
