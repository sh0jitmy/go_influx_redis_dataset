package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxDBの設定
const (
	token  = "your-influxdb-token"
	org    = "your-org"
	bucket = "your-bucket"
	url    = "http://localhost:8086" // InfluxDBのURL
)

// Redisの設定
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Redisサーバーのアドレス
})

// JSONデータの構造体
type Data struct {
	S        []int  `json:"s"`
	UpdateAt string `json:"updateat"`
}

// データを挿入する関数 (InfluxDB & Redis)
func insertData(client influxdb2.Client, data Data) error {
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// Time format の変換
	timestamp, err := time.Parse("2006-01-02 15:04:05.000", data.UpdateAt)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}

	// データをポイントとして追加 (InfluxDB用)
	p := influxdb2.NewPointWithMeasurement("signal_data").
		AddField("s", data.S).
		SetTime(timestamp)

	// InfluxDBにデータを書き込む
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		return err
	}

	// Redisにデータをセット (JSONとして保存)
	redisData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := redisClient.Set(context.Background(), "signal_data", redisData, 0).Err(); err != nil {
		return fmt.Errorf("failed to set Redis data: %v", err)
	}

	return nil
}

// 時間によるデータ検索関数 (InfluxDB)
func queryData(client influxdb2.Client, start, end time.Time) ([]api.Result, error) {
	queryAPI := client.QueryAPI(org)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => r._measurement == "signal_data")`, bucket, start.Format(time.RFC3339), end.Format(time.RFC3339))

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var results []api.Result
	for result.Next() {
		results = append(results, *result.Record())
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return results, nil
}

func main() {
	client := influxdb2.NewClient(url, token)
	defer client.Close()

	// 挿入データのサンプル
	data := Data{
		S:        []int{-127, -110, -100},
		UpdateAt: "2024-10-01 22:30:31.000",
	}

	// データの挿入 (InfluxDB & Redis)
	err := insertData(client, data)
	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		return
	}

	// 時間範囲のクエリ
	startTime, _ := time.Parse("2006-01-02T15:04:05Z", "2024-10-01T00:00:00Z")
	endTime := time.Now()

	results, err := queryData(client, startTime, endTime)
	if err != nil {
		fmt.Printf("Error querying data: %v\n", err)
		return
	}

	// 結果を出力
	for _, res := range results {
		fmt.Printf("Result: %+v\n", res)
	}
}
