package clock

import "time"

// SQL実行時などに利用する時刻情報を制御
// 基本的に永続化操作を行う際の時刻を固定化できるようにするのが目的
// tims.Time型はナノ秒単位の時刻制度の情報を持っており、
// 永続化したデータを取得して比較するとほぼ確実に時刻情報が不一致となるためである
// また、現在時刻の変化がテスト結果に影響することを回避する目的もある

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// FixedClocker はテスト用の固定時間を返す
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
