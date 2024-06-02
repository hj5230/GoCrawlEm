# GoCrawlEm

## Brief Intro

总之就是周末被亲戚抓住，要我帮ta爬一下国内主流媒体关于华为的报道，然后就写了这些，嗯。没有任何学习价值，都是最简单的实现，没必要写注释了。也没有代码质量一说，反正是一次性的能用就行... 主要用了chromedp：模拟chrome内核，bytedance/sonic：加载JSON（话说汇编是真的快，字节大佬nb）。至于为什么要传Github，只是用一下Codespace罢了... Codespace没有Chrome，跑一下Docker镜像。

In summary, over the weekend, a relative asked me to help with gathering news reports about Huawei from mainstream Chinese media outlets. I wrote these to accomplish that task. There is no learning value in this code as it is the simplest possible implementation and does not even require comments. Code quality was not a consideration since this is a one-time use script and only needs to be functional. The main libraries used were chromedp, which simulates a Chrome kernel, and bytedance/sonic for loading JSON (assembly is truly fast, the ByteDance engineers did an impressive job). As for why I uploaded this to GitHub, I was just need to utilize Codespace. Since Codespace does not have Chrome, remember to ran the Docker image.

## Usage

use pre-configured docker image:

```shell
docker run -d -p 9222:3000 --name chrome browserless/chrome
```
