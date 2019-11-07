package skeleton

import (
	"context"
	"log"

	"github.com/vilbert/go-skeleton/pkg/errors"

	"github.com/jmoiron/sqlx"

	skeletonEntity "github.com/vilbert/go-skeleton/internal/entity/skeleton"
)

type (
	// Data ...
	Data struct {
		db   *sqlx.DB
		stmt map[string]*sqlx.Stmt
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

const (
	getAllSkeletons  = "GetAllSkeletons"
	qGetAllSkeletons = "SELECT * FROM skeletons"

	getAllSkeletonsWithPage  = "GetAllSkeletonsWithPage"
	qGetAllSkeletonsWithPage = "SELECT * FROM skeletons LIMIT ?, ?"

	getSkeletonsCount  = "GetAllSkeletons"
	qGetSkeletonsCount = "SELECT COUNT(*) FROM skeletons"
)

var (
	readStmt = []statement{
		{getAllSkeletons, qGetAllSkeletons},
		{getAllSkeletonsWithPage, qGetAllSkeletonsWithPage},
		{getSkeletonsCount, qGetSkeletonsCount},
	}
)

// New ...
func New(db *sqlx.DB) Data {
	d := Data{
		db: db,
	}

	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// GetAllSkeletons ...
func (d Data) GetAllSkeletons(ctx context.Context) ([]skeletonEntity.Skeleton, error) {
	var (
		skeleton  skeletonEntity.Skeleton
		skeletons []skeletonEntity.Skeleton
		err       error
	)

	rows, err := d.stmt[getAllSkeletons].QueryxContext(ctx)

	for rows.Next() {
		if err := rows.StructScan(&skeleton); err != nil {
			return skeletons, errors.Wrap(err, "[DATA][GetAllSkeletons] ")
		}
		skeletons = append(skeletons, skeleton)
	}
	return skeletons, err
}

// GetAllSkeletonsWithPage ...
func (d Data) GetAllSkeletonsWithPage(ctx context.Context, offset int, length int) ([]skeletonEntity.Skeleton, error) {
	var (
		skeleton  skeletonEntity.Skeleton
		skeletons []skeletonEntity.Skeleton
		err       error
	)

	rows, err := d.stmt[getAllSkeletonsWithPage].QueryxContext(ctx, offset, length)

	for rows.Next() {
		if err := rows.StructScan(&skeleton); err != nil {
			return skeletons, errors.Wrap(err, "[DATA][GetAllSkeletonsWithPage] ")
		}
		skeletons = append(skeletons, skeleton)
	}
	return skeletons, err
}

// GetSkeletonsCount ...
func (d Data) GetSkeletonsCount(ctx context.Context) (int, error) {
	var (
		count int
		err   error
	)

	if err := d.stmt[getSkeletonsCount].QueryRowxContext(ctx).Scan(&count); err != nil {
		return count, errors.Wrap(err, "[DATA][GetSkeletonsCount] ")
	}
	return count, err

}
