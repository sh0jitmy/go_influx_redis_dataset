以下は、`README.md`に記載するためのコード解説と使い方です。

---

# InfluxDB & Redis Data Inserter with Time-based Query

このプロジェクトは、Go言語を使用して、JSON形式のデータをInfluxDBに挿入し、同時にRedisにデータを保存します。また、時間範囲に基づいてInfluxDBからデータをクエリする機能も提供します。

## 機能
- **InfluxDBへのデータ挿入**: JSON形式で受け取ったデータをInfluxDBに保存します。
- **Redisへのデータ保存**: InfluxDBにデータを挿入すると同時に、Redisにも同じデータを`SET`します。
- **時間範囲でのデータクエリ**: 指定した時間範囲内でInfluxDBからデータを検索できます。

## 前提条件

以下のツール・ライブラリが必要です:
- Go 1.19+
- InfluxDB (ローカルまたはクラウド)
- Redisサーバー (ローカルまたはクラウド)

## 使用方法

### 1. プロジェクトのセットアップ

まず、Goのプロジェクトを初期化します。

```bash
go mod init your_project
```

次に、必要なパッケージをインストールします。

```bash
go get github.com/influxdata/influxdb-client-go/v2
go get github.com/go-redis/redis/v8
```

### 2. 環境変数の設定

`main.go` 内の以下の設定をあなたの環境に合わせて変更してください。

```go
const (
    token  = "your-influxdb-token"
    org    = "your-org"
    bucket = "your-bucket"
    url    = "http://localhost:8086" // InfluxDBのURL
)

var redisClient = redis.NewClient(&redis.Options{
    Addr: "localhost:6379", // Redisサーバーのアドレス
})
```

- **InfluxDB Token**: InfluxDBの認証トークン。
- **Organization**: 使用するInfluxDBのオーガニゼーション名。
- **Bucket**: 保存先のバケット名。
- **InfluxDB URL**: InfluxDBのURL（ローカルまたはクラウド）。
- **Redis Addr**: Redisサーバーのアドレス。

### 3. データ挿入とクエリの実行

`main.go` 内でサンプルデータを挿入し、時間範囲でデータをクエリする処理が行われています。

#### データ挿入

```go
data := Data{
    S:        []int{-127, -110, -100},         // 信号データ
    UpdateAt: "2024-10-01 22:30:31.000",       // 更新日時
}
```

`insertData` 関数は、このデータをInfluxDBに挿入し、同時にRedisにも保存します。

#### データクエリ

指定した時間範囲でデータをクエリできます。

```go
startTime, _ := time.Parse("2006-01-02T15:04:05Z", "2024-10-01T00:00:00Z")
endTime := time.Now()
results, err := queryData(client, startTime, endTime)
```

クエリ結果はコンソールに出力されます。

### 4. 実行方法

以下のコマンドでプログラムを実行します。

```bash
go run main.go
```

データが正しく挿入され、RedisとInfluxDBに保存され、時間範囲に基づくクエリ結果が表示されます。

## 使用ライブラリ

- [InfluxDB Go Client](https://github.com/influxdata/influxdb-client-go) - InfluxDBと連携するためのクライアント。
- [Redis Go Client](https://github.com/go-redis/redis) - Redisと連携するためのクライアント。

## ライセンス

このプロジェクトはMITライセンスの下で提供されています。

---

これをプロジェクトのルートディレクトリに`README.md`として保存すれば、プロジェクトの使い方が簡単に説明できます。
