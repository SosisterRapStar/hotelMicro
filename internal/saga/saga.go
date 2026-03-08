package saga

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SosisterRapStar/LETI-paper/backoff"
	paperBroker "github.com/SosisterRapStar/LETI-paper/broker"
	"github.com/SosisterRapStar/LETI-paper/controller"
	"github.com/SosisterRapStar/LETI-paper/database"
	"github.com/SosisterRapStar/LETI-paper/message"
	"github.com/SosisterRapStar/LETI-paper/retry"
	"github.com/SosisterRapStar/LETI-paper/step"
	"github.com/google/uuid"
)

// flightReservedPayload — контракт сообщения топика flight.reserved (публикует flightsMicro).
// Совпадает с flightPayload в flightsMicro; hotel_booking_id добавляем в Execute и используем в Compensate.
type flightReservedPayload struct {
	BookingID       string `json:"booking_id"`
	UserID          string `json:"user_id"`
	FlightID        string `json:"flight_id"`
	HotelID         string `json:"hotel_id"`
	RoomID          string `json:"room_id"`
	CheckIn         string `json:"check_in"`
	CheckOut        string `json:"check_out"`
	AmountCents     int    `json:"amount_cents"`
	Currency        string `json:"currency"`
	FlightBookingID string `json:"flight_booking_id"`
	HotelBookingID  string `json:"hotel_booking_id"` // заполняем в Execute, нужен в Compensate
}

func parsePayload[T any](msg message.Message) (*T, error) {
	raw, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}
	return &result, nil
}

// Настройки ретраев для шагов саги.
//
// Поведение при ошибках:
//
// Execute:
//   - Ошибка без retry.AsRetryable(err) — ретраев нет, сразу вызывается OnError. OnError возвращает
//     кастомное сообщение и err → оно уходит в ErrorTopics (booking.created) с типом "failed" → flightsMicro
//     делает Compensate и шлёт flight.failed.
//   - Ошибка с retry.AsRetryable(err) — Execute повторяется до stepMaxRetries раз с экспоненциальным
//     бэкоффом (stepBackoffMin..stepBackoffMax). После исчерпания повторов вызывается OnError, далее как выше.
//   - OnError при успехе (return msg, nil) — сага продолжается, сообщение уходит в NextStepTopics.
//
// Compensate:
//   - Ошибка без AsRetryable — ретраев нет, вызывается OnCompensateError (если задан), иначе возвращается
//     ошибка компенсации.
//   - Ошибка с AsRetryable — Compensate повторяется до stepMaxRetries раз; после исчерпания — OnCompensateError
//     или возврат ошибки.
const (
	stepMaxRetries = 5
	stepBackoffMin = 100 * time.Millisecond
	stepBackoffMax = 3 * time.Second
)

type HotelSaga struct {
	Controller       *controller.Controller
	StepHotelReserve *step.Step
	ErrCh            chan error
}

func InitHotelSaga(
	ctx context.Context,
	db *sql.DB,
	pubsub paperBroker.Pubsub,
	tracing *controller.TracingConfig,
) (*HotelSaga, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}
	if pubsub == nil {
		return nil, fmt.Errorf("pubsub is required")
	}

	dbCtx := database.NewDBContext(db, database.SQLDialectMySQL)

	errCh := make(chan error, 128)

	ctrl, err := controller.New(&controller.Config{
		Subscriber: pubsub,
		Publisher:  pubsub,
		DB:         dbCtx,
		InfraRetry: &retry.Retrier{
			BackoffOptions: retry.BackoffOptions{
				BackoffPolicy: backoff.Expontential{},
				MinBackoff:    50 * time.Millisecond,
				MaxBackoff:    5 * time.Second,
			},
			MaxRetries: 10,
		},
		PollInterval:  1 * time.Second,
		BatchSize:     10,
		BackoffPolicy: backoff.Expontential{},
		BackoffMin:    100 * time.Millisecond,
		BackoffMax:    1 * time.Minute,
		ErrCh:         errCh,
		Tracing:       tracing,
	})
	if err != nil {
		return nil, fmt.Errorf("create controller: %w", err)
	}

	// Ретраи для Execute/Compensate: при retry.AsRetryable(err) шаг повторяется до stepMaxRetries раз.
	stepRetryPolicy := &retry.Retrier{
		BackoffOptions: retry.BackoffOptions{
			BackoffPolicy: backoff.Expontential{},
			MinBackoff:    stepBackoffMin,
			MaxBackoff:    stepBackoffMax,
		},
		MaxRetries: stepMaxRetries,
	}

	// При ошибке шлём в booking.created — этот топик читает flightsMicro (flight-reserve),
	// чтобы он выполнил Compensate и затем отправил flight.failed в booking.
	hotelReserveStep, err := step.New(&step.StepParams{
		Name:        "hotel-reserve",
		RetryPolicy: stepRetryPolicy,
		Routing: step.RoutingConfig{
			NextStepTopics: []string{TopicHotelReserved},
			ErrorTopics:    []string{TopicBookingCreated}, // слушает flightsMicro
		},
		// При падении Execute в ErrorTopics (booking.created) уходит кастомное сообщение
		// для flightsMicro: payload с flight_booking_id и остальными полями, чтобы flight-reserve.Compensate мог отменить бронь рейса.
		OnError: func(ctx context.Context, _ database.TxQueryer, msg message.Message, executeErr error) (message.Message, error) {
			p, err := parsePayload[flightReservedPayload](msg)
			if err != nil {
				return message.Message{}, fmt.Errorf("parse payload for error message: %w", err)
			}
			payload := map[string]any{
				"booking_id":        p.BookingID,
				"user_id":           p.UserID,
				"flight_id":         p.FlightID,
				"hotel_id":          p.HotelID,
				"room_id":           p.RoomID,
				"check_in":          p.CheckIn,
				"check_out":         p.CheckOut,
				"amount_cents":      p.AmountCents,
				"currency":          p.Currency,
				"flight_booking_id": p.FlightBookingID,
			}
			out := message.Message{}
			out.SagaID = msg.SagaID
			out.Payload = payload
			return out, executeErr
		},
		Execute: func(ctx context.Context, tx database.TxQueryer, msg message.Message) (message.Message, error) {
			p, err := parsePayload[flightReservedPayload](msg)
			if err != nil {
				return msg, err
			}
			if p.UserID == "" || p.HotelID == "" || p.RoomID == "" {
				return msg, fmt.Errorf("user_id, hotel_id and room_id are required")
			}

			checkIn, err := time.Parse(time.RFC3339, p.CheckIn)
			if err != nil {
				return msg, fmt.Errorf("parse check_in: %w", err)
			}
			checkOut, err := time.Parse(time.RFC3339, p.CheckOut)
			if err != nil {
				return msg, fmt.Errorf("parse check_out: %w", err)
			}

			bookingID := uuid.New().String()
			if _, err := tx.ExecContext(ctx, insertHotelBookingQuery,
				bookingID, p.UserID, p.HotelID, p.RoomID, checkIn, checkOut); err != nil {
				return msg, fmt.Errorf("insert hotel booking: %w", err)
			}

			result, err := tx.ExecContext(ctx, decrementRoomAvailableQuery, p.RoomID)
			if err != nil {
				return msg, fmt.Errorf("decrement rooms_available: %w", err)
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return msg, fmt.Errorf("rows affected: %w", err)
			}
			if affected == 0 {
				return msg, fmt.Errorf("no rooms available for room_id %s", p.RoomID)
			}

			if msg.Payload == nil {
				msg.Payload = make(map[string]any)
			}
			msg.Payload["hotel_booking_id"] = bookingID
			return msg, nil
		},
		Compensate: func(ctx context.Context, tx database.TxQueryer, msg message.Message) (message.Message, error) {
			p, err := parsePayload[flightReservedPayload](msg)
			if err != nil {
				return msg, err
			}
			if p.HotelBookingID == "" {
				return msg, fmt.Errorf("hotel_booking_id is required for compensation")
			}

			result, err := tx.ExecContext(ctx, cancelHotelBookingQuery, p.HotelBookingID)
			if err != nil {
				return msg, fmt.Errorf("cancel hotel booking: %w", err)
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return msg, fmt.Errorf("rows affected: %w", err)
			}
			if affected == 0 {
				return msg, nil
			}

			if _, err := tx.ExecContext(ctx, incrementRoomAvailableQuery, p.HotelBookingID); err != nil {
				return msg, fmt.Errorf("increment rooms_available: %w", err)
			}
			return msg, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create hotel-reserve step: %w", err)
	}

	return &HotelSaga{
		Controller:       ctrl,
		StepHotelReserve: hotelReserveStep,
		ErrCh:            errCh,
	}, nil
}
