# wecomrus
## Intro
This project can be embedded by logrus as a plugin for outputting the critical logs to `WeCom`. (`WeCom` is an instant-communication software that is similar as `WeChat` both developed by Tencent but focusing on daily working.)

## Installation
Install the package with go:
```go
go get github.com/henry0475/wecomrus
```

## Usage
1. For webhooks.
```go
import (
    "github.com/sirupsen/logrus"
    "github.com/henry0475/wecomrus"
)

func main() {
  log       := logrus.New()
  hook, err := wecomrus.NewWeComHook(
    wecomrus.Option{
        AppName: "Test APP",
    },
    wecomrus.Option{
        Webhooks: []string{
            "http://xxx/xxx/xxx/", // found from WeCom App creating a robot in a group chat
        }
    },
  )

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```

2. (NOT RECOMMENDED!) For group chat. It has less limitations but needs more privilege to work because you need to provide the corpID, corpSecret 
```go
import (
    "github.com/sirupsen/logrus"
    "github.com/henry0475/wecomrus"
)

func main() {
  log       := logrus.New()
  hook, err := wecomrus.NewWeComHook(
    wecomrus.Option{
        AppName: "Test APP",
    },
    wecomrus.Option{
        CorpID:       "Obtained from https://work.weixin.qq.com/",
        CorpSecret:   "Obtained from https://work.weixin.qq.com/",
        GroupChatID:  "xxxx", // You can get it from creating group chat phase by APIs.
    },
  )

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```
## Options
There are several options that can be used for customizing this plugin.
```go
wecomrus.Option{
    // AppName can be used for identifying the App
    AppName string
    // TimeFormat defines the format of time
    // Default: 01-02 15:04:05
    TimeFormat string
    // TimeZone defines the time location
    // Default: Asia/Chongqing
    // Notice: I recommend that you put `apk add tzdata` in the RUN statement for Debian
    TimeZone *time.Location
    // Safe defines whether the log message enables the safe flag or not
    // Default: SafeOff
    // ONLY valid for group chats, NOT for webhooks
    Safe SafeSwitcher
    // MessageFormart defines the different formats of the message
    // The {{app}} will be replaced by `AppName`
    // The {{time}} will be replaced by `TimeFormat`
    // The {{level}} will be replaced by different log's level
    // The {{message}} will be replaced based on the message the log has
    // The {{content}} will be replaced based on the log's fields info
    // `MessageFormater` is an interface that you should implement it for customizing it
    // Default: 
    // WeComRus Monitor
    // ***********
    // * AppName: {{app}}
    // * Time: {{time}}
    // * Level: {{level}}
    // ***********
    // * Message: {{message}}
    // ***********
    // {{content}}
    MessageFormart MessageFormater
    // MsgType defines what the specific type of log you want to display in WeCom
    // Default: TextMessage
    MsgType MessageType
    // Webhooks define a set of webhooks that can be used to deliver messages, also called robots
    // Because of the limitations from Tencent, each webhook is only allowed to send messages with the frequency of 20 messages per min. Therefore, I recommend that you assign multiple webhooks.
    // Currenly, all other requests over 20 message per min will be dropped, but you still can find your logs in your own system.
    Webhooks []string
    // EnableStats for opening the stats analysis.
    EnableStats *bool
    // GroupChatID defines the group ID in WeCom obtained from https://work.weixin.qq.com/
    GroupChatID string
    // CorpID defines the Corp ID obtained from https://work.weixin.qq.com/
    CorpID string
    // CorpSecret defines the Corp Secret obtained from https://work.weixin.qq.com/
    CorpSecret string
}
```

## Customized Output
1. Here is an example to customized the format of output message
```go
import (
    "github.com/sirupsen/logrus"
    "github.com/henry0475/wecomrus"
)

type myMessageFormat struct {
    content string
}

func (m myMessageFormat) ToString() string {
    return m.content
}

func main() {
    log := logrus.New()

    var myMessage myMessageFormat
    myMessage.content = `
{{app}}
Level: {{level}}
Msg: {{message}}
Time: {{time}}
----
This is a text message.`

    hook, err := wecomrus.NewWeComHook(
        wecomrus.Option{
            AppName: "Test APP",
        },
        wecomrus.Option{
            CorpID:       "Obtained from https://work.weixin.qq.com/",
            CorpSecret:   "Obtained from https://work.weixin.qq.com/",
            GroupChatID:  "xxxx",
        },
        wecomrus.Option{
            MessageFormart: myMessage,
        },
    )

    if err == nil {
        log.Hooks.Add(hook)
    }
}
```

## License
MIT