// 插入 interpolation

// 在你的项目中有支持过i18n吗？i18n.t() 模版引擎
// 拿i18next 来举例。通常情况下会将key和翻译维护在别的地方。比如这样：
// {
//   "evaluation": "BFE.dev is {{evaluation}}"
// }
// 然后在需要的地方，指定key以及需要的data即可得到想要的字符串。
// t('evaluation', {evaluation: 'fantastic'});
// "BFE.dev is fantastic"

function t(translation: string, data: any = {}): string {
  return translation.replace(/{{(.*?)}}/g, (_, g1) => data[g1] ?? '')
}

t('BFE.dev is {{{evaluation}', { evaluation: 'fantastic' })
// "BFE.dev is {{{evaluation}"

t('BFE.dev is {{{evaluation}}}', { '{evaluation': 'fantastic' })
// "BFE.dev is fantastic}"

t('BFE.dev is {{evalu ation}}', { 'evalu ation': 'fantastic' })
// "BFE.dev is fantastic"

t('BFE.dev is {{evaluation}}')
// "BFE.dev is "

export {}
