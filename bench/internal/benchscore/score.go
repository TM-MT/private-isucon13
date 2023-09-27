package benchscore

import (
	"context"
	"sync"

	"github.com/isucon/isucandar/score"
)

type ScoreTag string

const (
	SuccessGetTags score.ScoreTag = "success-get-tags"
	// ユーザ
	SuccessRegister     score.ScoreTag = "success-register"
	SuccessLogin        score.ScoreTag = "success-login"
	SuccessGetUser      score.ScoreTag = "success-get-user"
	SuccessGetUserTheme score.ScoreTag = "success-get-user-theme"
	// ライブ配信
	SuccessReserveLivestream  score.ScoreTag = "success-reserve-livestream"
	SuccessGetLivestreamByTag score.ScoreTag = "success-get-livestream-by-tag"
	// スパチャ
	SuccessGetSuperchats   score.ScoreTag = "success-get-superchats"
	SuccessPostSuperchat   score.ScoreTag = "success-post-superchat"
	SuccessReportSuperchat score.ScoreTag = "success-report-superchat"
	// リアクション
	SuccessGetReactions score.ScoreTag = "success-get-reactions"
	SuccessPostReaction score.ScoreTag = "success-post-reaction"
)

var (
	benchScore *score.Score
	// initOnce   sync.Once
	doneOnce sync.Once
)

func InitScore(ctx context.Context) {
	benchScore = score.NewScore(ctx)

	// FIXME: スコアの重み付けは後ほど考える
	benchScore.Set(SuccessGetTags, 1)
	benchScore.Set(SuccessRegister, 1)
	benchScore.Set(SuccessLogin, 1)
	benchScore.Set(SuccessGetUser, 1)
	benchScore.Set(SuccessGetUserTheme, 1)

	benchScore.Set(SuccessReserveLivestream, 1)
	benchScore.Set(SuccessGetLivestreamByTag, 1)

	benchScore.Set(SuccessGetSuperchats, 1)
	benchScore.Set(SuccessPostSuperchat, 1)
	benchScore.Set(SuccessReportSuperchat, 1)

	benchScore.Set(SuccessGetReactions, 1)
	benchScore.Set(SuccessPostReaction, 1)

	initProfit(ctx)
	initPenalty(ctx)
}

func AddScore(tag score.ScoreTag) {
	benchScore.Add(tag)
}

func GetCurrentScore() int64 {
	return benchScore.Sum()
}

func GetFinalScore() int64 {
	doneOnce.Do(func() {
		benchScore.Done()
	})
	return benchScore.Sum()
}
