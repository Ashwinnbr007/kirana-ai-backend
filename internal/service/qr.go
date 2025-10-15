package service

import (
	"context"
	"fmt"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
)

type QrService struct {
	db DatabasePort
}

func NewQrService(db DatabasePort) *QrService {
	return &QrService{
		db: db,
	}
}

func (q *QrService) ExtractContentsToDB(ctx context.Context, inventoryData models.InventoryData) error {

	err := q.db.WriteInventoryData(ctx, &inventoryData)
	if err != nil {
		return fmt.Errorf("could not write to db: %w", err)
	}

	return nil
}
