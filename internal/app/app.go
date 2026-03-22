package app

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	paperController "github.com/SosisterRapStar/cliros/controller"
	"github.com/SosisterRapStar/hotels/internal/adapter/controller"
	"github.com/SosisterRapStar/hotels/internal/adapter/controller/middleware"
	v1 "github.com/SosisterRapStar/hotels/internal/adapter/controller/v1"
	adapterKafka "github.com/SosisterRapStar/hotels/internal/adapter/kafka"
	adapterRepo "github.com/SosisterRapStar/hotels/internal/adapter/repository"
	"github.com/SosisterRapStar/hotels/internal/config"
	"github.com/SosisterRapStar/hotels/internal/domain/hotel"
	infrakafka "github.com/SosisterRapStar/hotels/internal/infrastructure/kafka"
	"github.com/SosisterRapStar/hotels/internal/infrastructure/telemetry"
	"github.com/SosisterRapStar/hotels/internal/saga"
)

type App struct {
	Controller *controller.Controller
}

func New(cfg *config.AppConfig) (*App, error) {
	rawDB, err := sql.Open("mysql", cfg.Repository.DSNMySQL())
	if err != nil {
		return nil, fmt.Errorf("opening mysql connection: %w", err)
	}
	rawDB.SetMaxIdleConns(cfg.Repository.MaxIdleConn)
	rawDB.SetMaxOpenConns(cfg.Repository.MaxOpenConn)
	rawDB.SetConnMaxIdleTime(cfg.Repository.MaxIdleLifetime)
	rawDB.SetConnMaxLifetime(cfg.Repository.MaxOpenLifetime)

	brokers := cfg.Kafka.Brokers
	if len(brokers) == 0 && cfg.Kafka.URL != "" {
		brokers = strings.Split(cfg.Kafka.URL, ",")
	}

	kafkaCfg := &infrakafka.Config{
		Brokers:          brokers,
		GroupID:          cfg.Kafka.GroupID,
		AckPolicy:        cfg.Kafka.Producer.AckPolicy,
		RetryMax:         cfg.Kafka.Producer.RetryMax,
		AutoCommitEnable: cfg.Kafka.Consumer.AutoCommitEnable,
		MaxWaitTime:      cfg.Kafka.Consumer.MaxWaitTime,
	}
	sagaPubsub, err := adapterKafka.NewSagaPubsub(kafkaCfg)
	if err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("create saga pubsub: %w", err)
	}

	if err := telemetry.Init(cfg); err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("init telemetry: %w", err)
	}

	ctx := context.Background()

	var sagaTracing *paperController.TracingConfig
	if cfg.Tracing.Enabled {
		sagaTracing = &paperController.TracingConfig{
			Tracer:     telemetry.Tracer("hotels-saga"),
			TracerName: "hotels",
		}
	}
	hotelSaga, err := saga.InitHotelSaga(ctx, rawDB, sagaPubsub, sagaTracing)
	if err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("init hotel saga: %w", err)
	}

	if err := hotelSaga.Controller.Register(saga.TopicFlightReserved, hotelSaga.StepHotelReserve); err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("register %s step: %w", saga.TopicFlightReserved, err)
	}

	if err := hotelSaga.Controller.Init(ctx); err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("init hotel saga controller: %w", err)
	}
	if err := sagaPubsub.Run(ctx); err != nil {
		_ = rawDB.Close()
		return nil, fmt.Errorf("run saga pubsub: %w", err)
	}

	sqlxDB := sqlx.NewDb(rawDB, "mysql")
	manager := adapterRepo.NewManager(sqlxDB)
	hotelRepo := adapterRepo.NewHotelRepository(sqlxDB, manager)
	hotelMod := hotel.NewModule(hotelRepo)

	return &App{
		Controller: &controller.Controller{
			Middleware: middleware.NewMiddleware(cfg),
			V1: v1.Controller{
				Dummy:   v1.NewDummyController(),
				Hotel:   v1.NewHotelController(hotelMod),
				Room:    v1.NewRoomController(hotelMod),
				Booking: v1.NewBookingController(hotelMod),
			},
		},
	}, nil
}
