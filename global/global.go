package global

import (
	"github.com/beastars1/lol-prophet-gui/conf"
	level "github.com/beastars1/lol-prophet-gui/pkg/logger"
	"log"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	UserInfo struct {
		IP string `json:"ip"`
	}
)

const (
	LogWriterCleanupKey = "logWriter"
	sentryDsn           = "https://815bc0eb2615452caa81e3ccb536ce05@o1184940.ingest.sentry.io/6303327"
	defaultLogPath      = "./logs/lol-prophet.log"
	AppName             = "LoL Prophet"
	Commit              = "dev"
	Version             = "1.4"
	ProjectUrl          = "https://github.com/beastars1/lol-prophet-gui"
)

const (
	horse1 = "汗血宝马"
	horse2 = "上等马"
	horse3 = "中等马"
	horse4 = "下等马"
	horse5 = "牛马"
)

var (
	horses = [5]conf.HorseScoreConf{
		{Score: 150, Name: horse1},
		{Score: 125, Name: horse2},
		{Score: 105, Name: horse3},
		{Score: 95, Name: horse4},
		{Score: 0.0001, Name: horse5},
	}
)

var (
	defaultShouldAutoOpenBrowserCfg = false
	DefaultClientConf               = conf.Client{
		AutoAcceptGame:                 false,
		AutoPickChampID:                0,
		AutoBanChampID:                 0,
		AutoSendTeamHorse:              false,
		ShouldSendSelfHorse:            true,
		HorseNameConf:                  [5]string{horse1, horse2, horse3, horse4, horse5},
		ChooseSendHorseMsg:             [5]bool{true, true, true, true, true},
		ChooseChampSendMsgDelaySec:     3,
		ShouldInGameSaveMsgToClipBoard: true,
		ShouldAutoOpenBrowser:          &defaultShouldAutoOpenBrowserCfg,
	}
	DefaultAppConf = conf.AppConf{
		Mode: conf.ModeProd,
		Sentry: conf.SentryConf{
			Enabled: true,
			Dsn:     sentryDsn,
		},
		PProf: conf.PProfConf{
			Enable: true,
		},
		Log: conf.LogConf{
			Level:    level.LevelInfoStr,
			Filepath: defaultLogPath,
		},
		CalcScore: conf.CalcScoreConf{
			Enabled:            true,
			FirstBlood:         [2]float64{10, 5},
			PentaKills:         [1]float64{20},
			QuadraKills:        [1]float64{10},
			TripleKills:        [1]float64{5},
			JoinTeamRateRank:   [4]float64{10, 5, 5, 10},
			GoldEarnedRank:     [4]float64{10, 5, 5, 10},
			HurtRank:           [2]float64{10, 5},
			Money2hurtRateRank: [2]float64{10, 5},
			VisionScoreRank:    [2]float64{10, 5},
			MinionsKilled: [][2]float64{
				{10, 20},
				{9, 10},
				{8, 5},
			},
			KillRate: []conf.RateItemConf{
				{Limit: 50, ScoreConf: [][2]float64{
					{15, 40},
					{10, 20},
					{5, 10},
				}},
				{Limit: 40, ScoreConf: [][2]float64{
					{15, 20},
					{10, 10},
					{5, 5},
				}},
			},
			HurtRate: []conf.RateItemConf{
				{Limit: 40, ScoreConf: [][2]float64{
					{15, 40},
					{10, 20},
					{5, 10},
				}},
				{Limit: 30, ScoreConf: [][2]float64{
					{15, 20},
					{10, 10},
					{5, 5},
				}},
			},
			AssistRate: []conf.RateItemConf{
				{Limit: 50, ScoreConf: [][2]float64{
					{20, 30},
					{18, 25},
					{15, 20},
					{10, 10},
					{5, 5},
				}},
				{Limit: 40, ScoreConf: [][2]float64{
					{20, 15},
					{15, 10},
					{10, 5},
					{5, 3},
				}},
			},
			AdjustKDA: [2]float64{2, 5},
			Horse:     horses,
			MergeMsg:  false,
		},
	}
	userInfo   = UserInfo{}
	confMu     = sync.Mutex{}
	Conf       = new(conf.AppConf)
	ClientConf = new(conf.Client)
	Logger     *zap.SugaredLogger
	Cleanups   = make(map[string]func() error)
)

// DB
var (
	SqliteDB *gorm.DB
)

func SetUserInfo(info UserInfo) {
	userInfo = info
}

func GetUserInfo() UserInfo {
	return userInfo
}

func Cleanup() {
	for name, cleanup := range Cleanups {
		if err := cleanup(); err != nil {
			log.Printf("%s cleanup err:%v\n", name, err)
		}
	}
	if fn, ok := Cleanups[LogWriterCleanupKey]; ok {
		_ = fn()
	}
}

func IsDevMode() bool {
	return GetEnv() == conf.ModeDebug
}

func GetEnv() conf.Mode {
	return Conf.Mode
}

func GetScoreConf() *conf.CalcScoreConf {
	confMu.Lock()
	defer confMu.Unlock()
	return &Conf.CalcScore
}

func SetScoreConf(scoreConf conf.CalcScoreConf) {
	confMu.Lock()
	Conf.CalcScore = scoreConf
	confMu.Unlock()
	return
}

func GetClientConf() *conf.Client {
	confMu.Lock()
	defer confMu.Unlock()
	data := *ClientConf
	return &data
}

func SetClientConf(cfg *conf.Client) *conf.Client {
	confMu.Lock()
	defer confMu.Unlock()
	if &cfg.AutoAcceptGame != nil {
		ClientConf.AutoAcceptGame = cfg.AutoAcceptGame
	}
	if &cfg.AutoPickChampID != nil {
		ClientConf.AutoPickChampID = cfg.AutoPickChampID
	}
	if &cfg.AutoBanChampID != nil {
		ClientConf.AutoBanChampID = cfg.AutoBanChampID
	}
	if &cfg.AutoSendTeamHorse != nil {
		ClientConf.AutoSendTeamHorse = cfg.AutoSendTeamHorse
	}
	if &cfg.ShouldSendSelfHorse != nil {
		ClientConf.ShouldSendSelfHorse = cfg.ShouldSendSelfHorse
	}
	if &cfg.HorseNameConf != nil {
		ClientConf.HorseNameConf = cfg.HorseNameConf
	}
	if &cfg.ChooseSendHorseMsg != nil {
		ClientConf.ChooseSendHorseMsg = cfg.ChooseSendHorseMsg
	}
	if &cfg.ChooseChampSendMsgDelaySec != nil {
		ClientConf.ChooseChampSendMsgDelaySec = cfg.ChooseChampSendMsgDelaySec
	}
	if &cfg.ShouldInGameSaveMsgToClipBoard != nil {
		ClientConf.ShouldInGameSaveMsgToClipBoard = cfg.ShouldInGameSaveMsgToClipBoard
	}
	if &cfg.ShouldAutoOpenBrowser != nil {
		ClientConf.ShouldAutoOpenBrowser = cfg.ShouldAutoOpenBrowser
	}
	return ClientConf
}
