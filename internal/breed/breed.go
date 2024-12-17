package breed

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(_ context.Context, db *pgxpool.Pool) (*Service, error) {
	return &Service{
		db: db,
	}, nil
}

func (s *Service) ListBreeds(ctx context.Context, in *BreedQuery) ([]*Breed, error) {
	return listBreeds(ctx, s.db, in)
}

type Breed struct {
	ID        string  `json:"id"`
	NameTH    string  `json:"nameTh"`
	NameEN    string  `json:"nameEn"`
	ShortName string  `json:"shortName"`
	Remark    *string `json:"remark"`
}

type BreedQuery struct {
	IDs        []string `json:"ids"`
	ShortNames []string `json:"shortNames"`
	Keyword    string   `json:"keyword"`
}

func (q *BreedQuery) toSql() (string, []interface{}, error) {
	eq := sq.Eq{}
	if len(q.IDs) > 0 {
		eq["id"] = q.IDs
	}
	if len(q.ShortNames) > 0 {
		eq["short_name"] = q.ShortNames
	}

	and := sq.And{eq}
	if q.Keyword != "" {
		and = append(and,
			sq.Expr(
				`name_th LIKE ? OR name_en LIKE ?`,
				fmt.Sprint("%", q.Keyword, "%"), fmt.Sprint("%", q.Keyword, "%"),
			),
		)
	}

	return and.ToSql()
}

func listBreeds(ctx context.Context, db *pgxpool.Pool, in *BreedQuery) ([]*Breed, error) {
	pred, args, err := in.toSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q, args := sq.
		Select(
			"id",
			"name_th",
			"name_en",
			"short_name",
			"remark",
		).
		From("breed").
		PlaceholderFormat(sq.Dollar).
		Where(pred, args...).
		MustSql()

	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	breeds := make([]*Breed, 0)
	for rows.Next() {
		var b Breed
		err := rows.Scan(
			&b.ID,
			&b.NameTH,
			&b.NameEN,
			&b.ShortName,
			&b.Remark,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		breeds = append(breeds, &b)
	}

	return breeds, nil
}
