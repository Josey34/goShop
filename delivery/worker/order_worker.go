package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Josey34/goshop/domain/entity"
)

func HandleOrderMessage(ctx context.Context, body string) error {
	var order entity.Order
	if err := json.Unmarshal([]byte(body), &order); err != nil {
		return err
	}

	log.Printf("[worker] processing order %s for customer %s", order.ID, order.CustomerID)
	return nil
}
