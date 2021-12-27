module github.com/taadis/blog-web

go 1.14

require (
	github.com/alicebob/miniredis/v2 v2.16.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.4
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/olivere/elastic/v7 v7.0.28
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/viper v1.7.1
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.1.3
	gorm.io/gorm v1.22.2
)

require (
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/micro/go-micro v1.18.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
