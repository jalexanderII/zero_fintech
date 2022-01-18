package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
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

type PaymentAction struct {
	AccountId       string    `json:"account_id,omitempty"`
	Amount          float32   `json:"amount,omitempty"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
	Status          int32     `json:"status,omitempty"`
}

type PaymentPlan struct {
	PaymentPlanId    string           `json:"payment_plan_id,omitempty"`
	UserId           string           `json:"user_id,omitempty"`
	PaymentTaskId    []string         `json:"payment_task_id,omitempty"`
	Timeline         float32          `json:"timeline,omitempty"`
	PaymentFreq      int32            `json:"payment_freq,omitempty"`
	AmountPerPayment float32          `json:"amount_per_payment,omitempty"`
	PlanType         int32            `json:"plan_type,omitempty"`
	EndDate          time.Time        `json:"end_date,omitempty"`
	Active           bool             `json:"active,omitempty"`
	Status           int32            `json:"status,omitempty"`
	PaymentAction    []*PaymentAction `json:"payment_action,omitempty"`
}

// CreateResponsePaymentPlan Takes in a model and returns a serializer
func CreateResponsePaymentPlan(paymentTaskModel *planning.PaymentPlan) PaymentPlan {
	paymentActions := make([]*PaymentAction, len(paymentTaskModel.GetPaymentAction()))
	for idx, paymentAction := range paymentTaskModel.GetPaymentAction() {
		paymentActions[idx] = &PaymentAction{
			AccountId:       paymentAction.AccountId,
			Amount:          paymentAction.Amount,
			TransactionDate: paymentAction.TransactionDate.AsTime(),
			Status:          int32(paymentAction.Status),
		}
	}
	return PaymentPlan{
		PaymentPlanId:    paymentTaskModel.PaymentPlanId,
		UserId:           paymentTaskModel.UserId,
		PaymentTaskId:    paymentTaskModel.PaymentTaskId,
		Timeline:         paymentTaskModel.Timeline,
		PaymentFreq:      int32(paymentTaskModel.PaymentFreq),
		AmountPerPayment: paymentTaskModel.AmountPerPayment,
		PlanType:         int32(paymentTaskModel.PlanType),
		EndDate:          paymentTaskModel.EndDate.AsTime(),
		Active:           paymentTaskModel.Active,
		Status:           int32(paymentTaskModel.Status),
		PaymentAction:    paymentActions,
	}
}

// CreateResponsePaymentTask Takes in a model and returns a serializer
func CreateResponsePaymentTask(paymentTaskModel *core.PaymentTask) PaymentTask {
	return PaymentTask{
		UserId:        paymentTaskModel.PaymentTaskId,
		TransactionId: paymentTaskModel.TransactionId,
		AccountId:     paymentTaskModel.AccountId,
		Amount:        paymentTaskModel.Amount,
		MetaData: MetaData{
			PreferredPlanType:    int32(paymentTaskModel.MetaData.PreferredPlanType),
			PreferredPaymentFreq: int32(paymentTaskModel.MetaData.PreferredPaymentFreq),
		}}
}

func GetPaymentPlan(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type GetPaymentPlanResponse struct {
			PaymentTasksIds []string `json:"payment_tasks_ids"`
		}
		var input GetPaymentPlanResponse
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
		paymentPlanResponse, err := client.GetPaymentPlan(ctx, &core.GetPaymentPlanRequest{PaymentTasksIds: input.PaymentTasksIds})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		responsePaymentPlans := make([]PaymentPlan, len(paymentPlanResponse.GetPaymentPlans()))
		for idx, paymentPlan := range paymentPlanResponse.GetPaymentPlans() {
			responsePaymentPlans[idx] = CreateResponsePaymentPlan(paymentPlan)
		}

		return c.Status(fiber.StatusOK).JSON(responsePaymentPlans)
	}
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
		return c.JSON(fiber.Map{"status": "success", "message": "Successfully created payment task", "payment_task": CreateResponsePaymentTask(paymentTask)})
	}
}

func ListPaymentTasks(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		listPaymentTasks, err := client.ListPaymentTasks(ctx, &core.ListPaymentTaskRequest{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		responsePaymentTasks := make([]PaymentTask, len(listPaymentTasks.GetPaymentTasks()))
		for idx, paymentTask := range listPaymentTasks.GetPaymentTasks() {
			responsePaymentTasks[idx] = CreateResponsePaymentTask(paymentTask)
		}

		return c.Status(fiber.StatusOK).JSON(responsePaymentTasks)
	}
}

func GetPaymentTask(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		getPaymentTask, err := client.GetPaymentTask(ctx, &core.GetPaymentTaskRequest{Id: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(CreateResponsePaymentTask(getPaymentTask))
	}
}

func UpdatePaymentTask(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type UpdatePaymentTaskResponse struct {
			Amount               float64 `json:"amount"`
			PreferredPlanType    int32   `json:"preferred_plan_type"`
			PreferredPaymentFreq int32   `json:"preferred_payment_freq"`
		}

		var updatePaymentTaskResponse UpdatePaymentTaskResponse
		if err := c.BodyParser(&updatePaymentTaskResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		updatePaymentTask, err := client.UpdatePaymentTask(ctx, &core.UpdatePaymentTaskRequest{
			Id: c.Params("id"),
			PaymentTask: &core.PaymentTask{
				Amount: updatePaymentTaskResponse.Amount,
				MetaData: &core.MetaData{
					PreferredPlanType:    common.PlanType(updatePaymentTaskResponse.PreferredPlanType),
					PreferredPaymentFreq: common.PaymentFrequency(updatePaymentTaskResponse.PreferredPaymentFreq),
				},
			},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(CreateResponsePaymentTask(updatePaymentTask))
	}
}

func DeletePaymentTask(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		response, err := client.DeletePaymentTask(ctx, &core.DeletePaymentTaskRequest{Id: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": response.GetStatus(), "data": CreateResponsePaymentTask(response.GetPaymentTask())})
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
			PreferredPlanType:    common.PlanType(paymentTask.MetaData.PreferredPlanType),
			PreferredPaymentFreq: common.PaymentFrequency(paymentTask.MetaData.PreferredPaymentFreq),
		},
	}
}
