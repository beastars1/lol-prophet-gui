package conf

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

const (
	SqliteDBPath = "prophet.db"
)

var (
	errBadConf = errors.New("错误的配置")
)

type (
	Client struct {
		AutoAcceptGame                 bool      `json:"autoAcceptGame"`                 // 自动接受
		AutoPickChampID                int       `json:"autoPickChampID"`                // 自动秒选
		AutoBanChampID                 int       `json:"autoBanChampID"`                 // 自动ban人
		AutoSendTeamHorse              bool      `json:"autoSendTeamHorse"`              // 是否自动发送消息到选人界面
		ShouldSendSelfHorse            bool      `json:"shouldSendSelfHorse"`            // 是否发送自己马匹信息
		HorseNameConf                  [5]string `json:"horseNameConf"`                  // 马匹名称配置
		ChooseSendHorseMsg             [5]bool   `json:"chooseSendHorseMsg"`             // 选择发送哪些马匹信息
		ChooseChampSendMsgDelaySec     int       `json:"chooseChampSendMsgDelaySec"`     // 选人阶段延迟几秒发送
		ShouldInGameSaveMsgToClipBoard bool      `json:"shouldInGameSaveMsgToClipBoard"` // 进入对局后保存敌方马匹消息到剪切板中
		ShouldAutoOpenBrowser          *bool     `json:"shouldAutoOpenBrowser"`          // 是否自动打开浏览器
	}
)

func ValidClientConf(conf *Client) error {
	for _, s := range conf.HorseNameConf {
		if s == "" {
			return errBadConf
		}
	}
	return nil
}

func (conf *Client) Update() error {
	body, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	_, err = http.Post("http://127.0.0.1:4396/v1/config/update", "", bytes.NewBuffer(body))
	return err
}
