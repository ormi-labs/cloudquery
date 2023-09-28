package client

import (
	"context"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/jackc/pgx/v5/pgxpool"
)

type query struct {
	entity EntityName
	arg1   string
}

type queryPool struct {
	entities map[EntityName]*entity
	conn     *pgxpool.Pool
	in       chan query
	out      chan error
	quit     chan struct{}
	res      chan<- message.SyncMessage
}

func (p *queryPool) startWorkers(ctx context.Context) {
	for i := 0; i < len(p.entities); i++ {
		go func(i int) {
		NEXT:
			for q := range p.in {
				e := p.entities[q.entity]
				rows, err := p.conn.Query(ctx, e.sqlStmt, q.arg1)
				if err != nil {
					p.out <- err
					continue NEXT
				}

				for rows.Next() {
					values, err := rows.Values()
					if err != nil {
						rows.Close()
						p.out <- err
						continue NEXT
					}

					for i, value := range values {
						val, err := e.transformers[i](value)
						if err != nil {
							rows.Close()
							p.out <- err
							continue NEXT
						}

						s := scalar.NewScalar(e.arrowSchema.Field(i).Type)
						if err := s.Set(val); err != nil {
							rows.Close()
							p.out <- err
							continue NEXT
						}

						scalar.AppendToBuilder(e.recBuilder.Field(i), s)
					}
				}
				rows.Close()
				p.res <- &message.SyncInsert{Record: e.recBuilder.NewRecord()}
				p.out <- nil
			}
		}(i)
	}
}

func (p *queryPool) start(ctx context.Context) {
	p.startWorkers(ctx)
}

func (p *queryPool) stop() {
	close(p.quit)
	close(p.in) // we know that stop is called last
}

func (p *queryPool) add(q query) {
	select {
	case p.in <- q:
	case <-p.quit:
	}
}
