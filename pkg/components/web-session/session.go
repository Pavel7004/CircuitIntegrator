package websession

import (
	"context"
	"errors"

	"github.com/Pavel7004/Common/tracing"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	pointgenerator "github.com/Pavel7004/GraphPlot/pkg/components/point-generator"
	"github.com/Pavel7004/GraphPlot/pkg/domain"
)

var (
	ErrInterrupted = errors.New("Integration was interrupted by user")
)

type Session struct {
	conn *websocket.Conn
}

func New(conn *websocket.Conn) *Session {
	s := new(Session)

	s.conn = conn

	go s.listenConn()

	return s
}

func (s *Session) listenConn() {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	var (
		data = &domain.CircuitQuery{}

		endCh chan struct{}
	)
	for {
		if err := s.conn.ReadJSON(data); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Error().Err(err).Msg("Error on reading json from user")
			}
			break
		}
		if endCh != nil {
			endCh <- struct{}{}

			close(endCh)
		}
		if err := data.Check(); err != nil {
			log.Warn().Err(err).Msg("User data is incorrect")
		}

		circuit := circuit.New(&circuit.ChargeComponents{
			SupplyVoltage:     data.SupplyVoltage,
			Capacity:          data.Capacity,
			Resistance:        data.Resistance,
			StagesCount:       data.StagesCount,
			GapTriggerVoltage: data.GapTriggerVoltage,
			HoldingVoltage:    data.HoldingVoltage,
		}, &circuit.LoadComponents{
			Resistance: data.LoadResistance,
		})

		endCh = make(chan struct{})
		go s.plot(ctx, endCh, data.IntNum, circuit, data.Step)
	}
}

func (s *Session) plot(ctx context.Context, endCh chan struct{}, intNum int, circ *circuit.Circuit, step float64) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	bufferX := make([]float64, 0, 10)
	bufferY := make([]float64, 0, 10)

	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: circ,
		Step:    step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			select {
			case <-endCh:
				return ErrInterrupted
			default:
				if len(bufferX) != cap(bufferX) {
					bufferX = append(bufferX, t)
					bufferY = append(bufferY, x.GetLoadVoltage())

					bufferX = make([]float64, 0, 10)
					bufferY = make([]float64, 0, 10)
				} else if err := s.conn.WriteJSON(struct {
					Type string    `json:"type"`
					X    []float64 `json:"x"`
					Y    []float64 `json:"y"`
				}{
					"point",
					bufferX,
					bufferY,
				}); err != nil {
					return err
				}
			}
			return nil
		},
		NewIntFn: domain.Integrators[intNum],
	})

	if err := s.conn.WriteJSON(struct {
		Type string `json:"type"`
	}{
		"end",
	}); err != nil {
		log.Error().Err(err).Msg("Failed to send \"end\" signal")
	}

}
