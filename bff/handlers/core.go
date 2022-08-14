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
	PreferredPlanType         int32   `json:"preferred_plan_type"`
	PreferredTimelineInMonths float64 `json:"preferred_timeline_in_months"`
	PreferredPaymentFreq      int32   `json:"preferred_payment_freq"`
}

type PaymentTask struct {
	UserId    string  `json:"user_id"`
	AccountId string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

type PaymentAction struct {
	AccountId       string    `json:"account_id,omitempty"`
	Amount          float64   `json:"amount,omitempty"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
	Status          int32     `json:"status,omitempty"`
}

type PaymentPlan struct {
	PaymentPlanId    string           `json:"payment_plan_id,omitempty"`
	UserId           string           `json:"user_id,omitempty"`
	PaymentTaskId    []string         `json:"payment_task_id,omitempty"`
	Timeline         float64          `json:"timeline,omitempty"`
	PaymentFreq      int32            `json:"payment_freq,omitempty"`
	Amount           float64          `json:"amount,omitempty"`
	AmountPerPayment float64          `json:"amount_per_payment,omitempty"`
	PlanType         int32            `json:"plan_type,omitempty"`
	EndDate          time.Time        `json:"end_date,omitempty"`
	Active           bool             `json:"active,omitempty"`
	Status           int32            `json:"status,omitempty"`
	PaymentAction    []*PaymentAction `json:"payment_action,omitempty"`
}

type AccountInfo struct {
	TransactionIds []string `json:"transaction_ids,omitempty"`
	AccountId      string   `json:"account_id,omitempty"`
	Amount         float64  `json:"amount,omitempty"`
}

type GetPaymentPlanRequest struct {
	AccountInfo []AccountInfo `json:"account_info,omitempty"`
	UserId      string        `json:"user_id,omitempty"`
	MetaData    MetaData      `json:"meta_data,omitempty"`
	SavePlan    bool          `json:"save_plan"`
}

// CreateResponsePaymentPlan Takes in a model and returns a serializer
func CreateResponsePaymentPlan(paymentTaskModel *common.PaymentPlan) PaymentPlan {
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
		Amount:           paymentTaskModel.Amount,
		AmountPerPayment: paymentTaskModel.AmountPerPayment,
		PlanType:         int32(paymentTaskModel.PlanType),
		EndDate:          paymentTaskModel.EndDate.AsTime(),
		Active:           paymentTaskModel.Active,
		Status:           int32(paymentTaskModel.Status),
		PaymentAction:    paymentActions,
	}
}

// CreateResponsePaymentTask Takes in a model and returns a serializer
func CreateResponsePaymentTask(paymentTaskModel *common.PaymentTask) PaymentTask {
	return PaymentTask{
		UserId:    paymentTaskModel.PaymentTaskId,
		AccountId: paymentTaskModel.AccountId,
		Amount:    paymentTaskModel.Amount,
	}
}

func GetPaymentPlan(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var input GetPaymentPlanRequest
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
		accountInfoList := make([]*core.AccountInfo, len(input.AccountInfo))
		for idx, accountInfo := range input.AccountInfo {
			accountInfoList[idx] = AccountInfoDBToPB(accountInfo)
		}
		metaData := &common.MetaData{
			PreferredPlanType:         common.PlanType(input.MetaData.PreferredPlanType),
			PreferredTimelineInMonths: input.MetaData.PreferredTimelineInMonths,
			PreferredPaymentFreq:      common.PaymentFrequency(input.MetaData.PreferredPaymentFreq),
		}
		paymentPlanResponse, err := client.GetPaymentPlan(ctx, &core.GetPaymentPlanRequest{AccountInfo: accountInfoList, UserId: input.UserId, MetaData: metaData, SavePlan: input.SavePlan})
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
		paymentTask, err := client.CreatePaymentTask(ctx, &common.CreatePaymentTaskRequest{PaymentTask: PaymentTaskDBToPB(input)})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create payment task", "data": err})
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Successfully created payment task", "payment_task": CreateResponsePaymentTask(paymentTask)})
	}
}

func ListPaymentTasks(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		listPaymentTasks, err := client.ListPaymentTasks(ctx, &common.ListPaymentTaskRequest{})
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
		getPaymentTask, err := client.GetPaymentTask(ctx, &common.GetPaymentTaskRequest{Id: c.Params("id")})
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

		updatePaymentTask, err := client.UpdatePaymentTask(ctx, &common.UpdatePaymentTaskRequest{
			Id: c.Params("id"),
			PaymentTask: &common.PaymentTask{
				Amount: updatePaymentTaskResponse.Amount,
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
		response, err := client.DeletePaymentTask(ctx, &common.DeletePaymentTaskRequest{Id: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": response.GetStatus(), "data": CreateResponsePaymentTask(response.GetPaymentTask())})
	}
}

func GetUserAccounts(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accounts, err := client.ListUserAccounts(ctx, &core.ListUserAccountsRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on fetching user's accounts", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": accounts})
	}
}

func GetUserDebitAccountBalance(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		debitAccBalance, err := client.GetDebitAccountBalance(ctx, &core.GetDebitAccountBalanceRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on fetching user's debit account", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": debitAccBalance})
	}
}

func GetUserTransactions(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		transactions, err := client.ListUserTransactions(ctx, &core.ListUserTransactionsRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on fetching user's transactions", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": transactions})
	}
}

func GetUserPaymentPlans(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		plans, err := client.ListUserPaymentPlans(ctx, &common.ListUserPaymentPlansRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error on fetching user's transactions", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": plans})
	}
}

func GetWaterfallOverview(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		overview, err := client.GetWaterfallOverview(ctx, &planning.GetUserOverviewRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error fetching user's waterfall", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": overview})
	}
}

func GetAmountPaidPercentage(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		percentage, err := client.GetAmountPaidPercentage(ctx, &planning.GetUserOverviewRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error fetching user's amounts paid", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": percentage})
	}
}

func GetPercentageCoveredByPlans(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		plans, err := client.GetPercentageCoveredByPlans(ctx, &planning.GetUserOverviewRequest{UserId: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error fetching user's percent plan covered", "data": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": plans})
	}
}

// PaymentTaskDBToPB converts a PaymentTask DB object to its proto object
func PaymentTaskDBToPB(paymentTask PaymentTask) *common.PaymentTask {
	return &common.PaymentTask{
		UserId:    paymentTask.UserId,
		AccountId: paymentTask.AccountId,
		Amount:    paymentTask.Amount,
	}
}

// AccountInfoDBToPB converts a AccountInfo DB object to its proto object
func AccountInfoDBToPB(accountInfo AccountInfo) *core.AccountInfo {
	return &core.AccountInfo{
		TransactionIds: accountInfo.TransactionIds,
		AccountId:      accountInfo.AccountId,
		Amount:         accountInfo.Amount,
	}
}
