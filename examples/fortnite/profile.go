package fortnite

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/usk81/bone"
)

// BRPlayerStats defines the response of Fornite BR Player Stats API
type BRPlayerStats struct {
	AccountID        string         `json:"accountId"`
	PlatformID       int            `json:"platformId"`
	PlatformName     string         `json:"platformName"`
	PlatformNameLong string         `json:"platformNameLong"`
	EpicUserHandle   string         `json:"epicUserHandle"`
	Stats            Stats          `json:"stats"`
	LifeTimeStats    []LifeTimeStat `json:"lifeTimeStats"`
	RecentMatches    []Match        `json:"recentMatches"`
}

// Stats is Stats data in BRPlayerStats
type Stats struct {
	Solo                  StatField `json:"p2"`
	Duo                   StatField `json:"p10"`
	Squad                 StatField `json:"p9"`
	LimitTimeModes        StatField `json:"ltm"`
	CurrentSolo           StatField `json:"curr_p2"`
	CurrentDuo            StatField `json:"curr_p10"`
	CurrentSquad          StatField `json:"curr_p9"`
	CurrentLimitTimeModes StatField `json:"curr_ltm"`
}

// StatField is a part of Stats
type StatField struct {
	TrnRating     ScoreField `json:"trnRating"`
	Score         ScoreField `json:"score"`
	Top1          ScoreField `json:"top1"`
	Top3          ScoreField `json:"top3"`
	Top5          ScoreField `json:"top5"`
	Top6          ScoreField `json:"top6"`
	Top10         ScoreField `json:"top10"`
	Top12         ScoreField `json:"top12"`
	Top25         ScoreField `json:"top25"`
	KD            RatioField `json:"kd"`
	WinRatio      RatioField `json:"winRatio"`
	Matches       ScoreField `json:"matches"`
	Kills         ScoreField `json:"kills"`
	MinutesPlayed ScoreField `json:"minutesPlayed"`
	KPM           RatioField `json:"kpm"`
	KPG           RatioField `json:"kpg"`
	AvgTimePlayed RatioField `json:"avgTimePlayed"`
	ScorePerMatch RatioField `json:"scorePerMatch"`
	ScorePerMin   RatioField `json:"scorePerMin"`
}

// ScoreField stores personal score
type ScoreField struct {
	Label        string  `json:"label"`
	Field        string  `json:"field"`
	Category     string  `json:"category"`
	ValueInt     int     `json:"valueInt"`
	Value        string  `json:"value"`
	Rank         int     `json:"rank"`
	Percentile   float64 `json:"percentile"`
	DisplayValue string  `json:"displayValue"`
}

// RatioField stores personal ratio
type RatioField struct {
	Label        string  `json:"label"`
	Field        string  `json:"field"`
	Category     string  `json:"category"`
	ValueDec     float64 `json:"valueDec"`
	Value        string  `json:"value"`
	Rank         int     `json:"rank"`
	Percentile   float64 `json:"percentile"`
	DisplayValue string  `json:"displayValue"`
}

// LifeTimeStat stores the key and value for lifetime stats map
type LifeTimeStat struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Match defines a part of recent matches
type Match struct {
	ID              int     `json:"id"`
	AccountID       string  `json:"accountId,omitempty"`
	Top1            int     `json:"top1"`
	Top3            int     `json:"top3"`
	Top5            int     `json:"top5"`
	Top6            int     `json:"top6"`
	Top10           int     `json:"top10"`
	Top12           int     `json:"top12"`
	Top25           int     `json:"top25"`
	DateCollected   string  `json:"dateCollected"`
	Kills           int     `json:"kills"`
	Matches         int     `json:"matches"`
	MinutesPlayed   int     `json:"minutesPlayed,omitempty"`
	Platform        int     `json:"platform,omitempty"`
	PlayersOutlived int     `json:"playersOutlived"`
	Playlist        string  `json:"playlist"`
	PlaylistID      int     `json:"playlistId,omitempty"`
	Score           int     `json:"score"`
	TrnRating       float64 `json:"trnRating,omitempty"`
	TrnRatingChange float64 `json:"trnRatingChange,omitempty"`
}

type ProfileService struct {
	client bone.Client
}

func (s *ProfileService) SetClient(c bone.Client) {
	s.client = c
}

func (s *ProfileService) Stats(platform, nickname string) (result BRPlayerStats, err error) {
	c, ok := s.client.(*bone.DefaultClient)
	if !ok {
		err = errors.New("client is invalid")
		return
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("profile/%s/%s", platform, nickname), nil, nil)
	if err != nil {
		return
	}
	if _, err = c.Do(nil, req, bone.JSONDecode, &result); err != nil {
		return
	}
	return result, nil
}

func (s *ProfileService) MatchHistory(accountID string) (result []Match, err error) {
	c, ok := s.client.(*bone.DefaultClient)
	if !ok {
		err = errors.New("client is invalid")
		return
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("profile/account/%s/matches", accountID), nil, nil)
	if err != nil {
		return
	}
	if _, err = c.Do(nil, req, bone.JSONDecode, &result); err != nil {
		return
	}
	return result, nil
}
