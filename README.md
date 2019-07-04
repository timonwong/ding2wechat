# ding2wechat

一个模拟钉钉 webhook 服务器，然后转发到企业微信机器人的转发器。这样无需修改代码，就可以将企业微信机器人能力，接入到支持发送钉钉机器人消息的应用中。

## 使用

参考配置文件:

```yaml
receivers:
  - name: without_mention
    targets:
      - url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

  - name: mention_list
    targets:
      - url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        mentioned_list: ["bot", "@all"]  # (text only) userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list

  - name: mention_mobile
    targets:
      - url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
        mentioned_mobile_list: ["13800001111", "@all"] # (text only) 手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
```
