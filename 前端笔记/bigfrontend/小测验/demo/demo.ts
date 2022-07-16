import { crawl } from '../../../crawl'

crawl({
  url: 'https://bigfrontend.dev/zh/quiz?',
  xPathExpression: "//ul[@class='List__ListItems-sc-1p5i700-1 kUISA-D']//li//text()",
  extName: '.ts',
})
