package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
)

// MetaData is a DB Serialization of Proto MetaData
type MetaData struct {
	PreferredPlanType    int32 `json:"preferred_plan_type"`
	PreferredPaymentFreq int32 `json:"preferred_payment_freq"`
}

type PaymentTask struct {
	UserId        string   `json:"user_id"`
	TransactionId string   `json:"transaction_id"`
	AccountId     string   `json:"account_id"`
	Amount        float64  `json:"amount"`
	MetaData      MetaData `json:"meta_data"`
}

func CreatePaymentTask(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var input PaymentTask

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
		paymentTask, err := client.CreatePaymentTask(ctx, &core.CreatePaymentTaskRequest{PaymentTask: PaymentTaskDBToPB(input)})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create payment task", "data": err})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Successfully created payment task", "payment_task": paymentTask})
	}
}

// PaymentTaskDBToPB converts a PaymentTask DB object to its proto object
func PaymentTaskDBToPB(paymentTask PaymentTask) *core.PaymentTask {
	return &core.PaymentTask{
		UserId:        paymentTask.UserId,
		TransactionId: paymentTask.AccountId,
		AccountId:     paymentTask.TransactionId,
		Amount:        paymentTask.Amount,
		MetaData: &core.MetaData{
			PreferredPlanType:    core.PlanType(paymentTask.MetaData.PreferredPlanType),
			PreferredPaymentFreq: core.PaymentFrequency(paymentTask.MetaData.PreferredPaymentFreq),
		},
	}
}
