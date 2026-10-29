package main

import (
	"time"

	"github.com/Raghav847/social-go/internal/db"
	"github.com/Raghav847/social-go/internal/env"
	"github.com/Raghav847/social-go/internal/mailer"
	"github.com/Raghav847/social-go/internal/store"
	"go.uber.org/zap"
)

//	@title			social-go API
//	@description	API for social-go, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			//fromEmail: env.Ge....  hello@demomailtrap.co
			exp:       time.Hour * 24 * 3,
			fromEmail: "hello@demomailtrap.co",
			mailtrap: mailtrapConfig{
				apiKey: env.GetString("MAILTRAP_API_TOKEN", ""),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("db connection established")

	mailtrap, err := mailer.NewMailtrapClient(cfg.mail.mailtrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	store := store.NewStore(db)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailtrap,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
