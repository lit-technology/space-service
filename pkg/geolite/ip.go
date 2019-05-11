package geolite

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/philip-bui/space-service/services/postgres"
	"github.com/rs/zerolog/log"
)

// IPBatch is a transaction handling start, batching of IP and commit of the transaction.
type IPBatch struct {
	*sql.Tx
	*sql.Stmt
}

// BeginIPBatch begins a IP Batch Transaction.
func BeginIPBatch() (*IPBatch, error) {
	txn, err := postgres.DB.Begin()
	if err != nil {
		log.Error().Err(err).Msg("error beginning transaction")
		return nil, err
	}
	stmt, err := txn.Prepare(pq.CopyIn(postgres.TableIP, postgres.ColIP, postgres.ColGeoliteID, postgres.ColCellID))
	if err != nil {
		log.Error().Err(err).Msg("error preparing statement")
		return nil, err
	}
	return &IPBatch{
		txn,
		stmt,
	}, nil
}

// AddRow adds a IP row with ip, geoliteID and cellID.
func (g *IPBatch) AddRow(ip string, geoliteID int, cellID int64) error {
	if _, err := g.Stmt.Exec(ip, geoliteID, cellID); err != nil {
		log.Error().Err(err).Str("ip", ip).Int("geoliteID", geoliteID).Int64("cellID", cellID).Msg("error adding row")
		return err
	}
	return nil
}

// ExecAndCommit executes and commits the batch transaction.
func (g *IPBatch) ExecAndCommit() error {
	if _, err := g.Stmt.Exec(); err != nil {
		log.Error().Err(err).Msg("error executing statement")
		return err
	}
	if err := g.Close(); err != nil {
		log.Error().Err(err).Msg("error closing statement")
		return err
	}
	if err := g.Commit(); err != nil {
		log.Error().Err(err).Msg("error commiting transaction")
		return err
	}
	return nil
}
