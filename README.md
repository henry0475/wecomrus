# wecomrus
## Intro
This project can be embedded by logrus as a plugin for outputting the critical logs to `WeCom`. (`WeCom` is an instant-communication software that is similar as `WeChat` both developed by Tencent but focusing on daily working.)

## Installation
Install the package with go:
```go
go get github.com/henry0475/wecomrus
```

## Usage
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
        GroupChatID:  "xxxx",
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
    // Notice: If your applciation will be running on the Docker platform with alpine:latest, the plugin will **throw panic**
    // As long as you put `apk add tzdata` in the RUN statement for Debian
    TimeZone *time.Location
    // Safe defines whether the log message enables the safe flag or not
    // Default: SafeOff
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

    // --- You MUST assign the following options! --- //

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