package main

import (
	"math"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Data struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IsFireDetected bool   `json:"isFireDetected"`
	IsDroneGoingTo bool   `json:"isDroneGoingTo"`
}

func main() {
	// ローカルでの動作確認用のサーバーを起動
	// isRunWithMockServer := false
	// if isRunWithMockServer {
	// 	runIsDroneStartServer()
	// }

	drone := tello.NewDriver("8889")

	work := func() {
		drone.TakeOff()

		// 3秒間前進
		gobot.After(3*time.Second, func() {
			drone.Forward(50)
		})

		// 前進後に2秒間ホバリング
		gobot.After(6*time.Second, func() {
			drone.Hover()
		})

		// 半径50cmの範囲を4周回 その後ホバリングへ移行
		gobot.After(10*time.Second, func() {
			for j := 0; j < 4; j++ { // 4回周回するためのループ
				for i := 0; i < 360; i += 10 {
					rad := float64(i) * (math.Pi / 180)
					x := 50 * math.Cos(rad)
					y := 50 * math.Sin(rad)
					drone.Forward(int(x))
					drone.Right(int(y))
					time.Sleep(100 * time.Millisecond)
					drone.Hover()
				}
			}
		})

		// 28秒後に後退開始(3秒間後退)
		gobot.After(28*time.Second, func() {
			drone.Backward(50)
		})

		// 後退後に2秒間ホバリング
		gobot.After(31*time.Second, func() {
			drone.Hover()
		})
		// 着陸
		gobot.After(33*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Connections().Start()

	defer func() {
		robot.Stop()
		os.Exit(0)
	}()

	time.Sleep(1 * time.Second)

	robot.Start()

}
