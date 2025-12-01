package main

import (
	"time"

	"github.com/ecetinerdem/forseerv2/internal/auth"
	"github.com/ecetinerdem/forseerv2/internal/db"
	"github.com/ecetinerdem/forseerv2/internal/env"
	"github.com/ecetinerdem/forseerv2/internal/mailer"
	"github.com/ecetinerdem/forseerv2/internal/store"
	"github.com/ecetinerdem/forseerv2/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

//	@title			ForSeer API
//	@description	A platform for analyze your portfolio with power of AI.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.forseer.com
//	@contact.email	forseerbussiness@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		frontEndURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/forseer?sslmode=disable"),
			maxOpenConn: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdleConn: env.GetInt("DB_MAX_IDLE_CONN", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "1.0.0"),
		mail: mailConfig{
			sendGrid: sendGridConfig{
				apiKey: env.GetString("MAILER_API_KEY", ""),
			},
			fromEmail: env.GetString("MAILER_FROM_EMAIL", ""),
			expiry:    time.Hour * 24 * 3,
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", ""),
				pass: env.GetString("AUTH_BASIC_PASS", ""),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", ""),
				expiry: time.Hour * 24 * 3,
				iss:    env.GetString("AUTH_TOKEN_ISS", "forseer"),
			},
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// DB
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConn, cfg.db.maxIdleConn, cfg.db.maxIdleTime)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection established")

	//Store
	store := store.NewStorage(db)

	//Mailer
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	//JWT Authenticator
	JWTAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	//Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis connection established")

	}

	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: JWTAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
