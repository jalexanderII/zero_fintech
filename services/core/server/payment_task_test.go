package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreServer_CreatePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	equinoxPT := &core.PaymentTask{
		UserId:        "61df93c0ac601d1be8e64613",
		TransactionId: "61dfa20adebb9d4fb62b9703",
		AccountId:     "61df9b621d2c2b15a6e53ec9",
		Amount:        325,
		MetaData: &core.MetaData{
			PreferredPlanType:    core.PlanType_PLANTYPE_OPTIM_CREDIT_SCORE,
			PreferredPaymentFreq: core.PaymentFrequency_PAYMENTFREQ_MONTHLY,
		},
	}

	paymentTask, err := server.CreatePaymentTask(ctx, &core.CreatePaymentTaskRequest{PaymentTask: equinoxPT})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
	equinoxPT.PaymentTaskId = paymentTask.PaymentTaskId
	if paymentTask != equinoxPT {
		t.Errorf("2: Error creating account with correct details: %v", paymentTask)
	}
}

func TestCoreServer_GetPaymentTask(t *testing.T) {
	server, ctx := GenServer()

	paymentTask, err := server.GetPaymentTask(ctx, &core.GetPaymentTaskRequest{Id: "61dfa8296c734067e6726761"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if paymentTask.Amount != 325 {
		t.Errorf("2: Failed to fetch correct paymentTask: %+v", paymentTask)
	}
}

func TestCoreServer_ListPaymentTasks(t *testing.T) {
	server, ctx := GenServer()

	paymentTasks, err := server.ListPaymentTasks(ctx, &core.ListPaymentTaskRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(paymentTasks.PaymentTasks) < 1 {
		t.Errorf("2: Failed to fetch realtors: %+v", paymentTasks.PaymentTasks[0])
	}
}

func TestCoreServer_UpdatePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	u := &core.PaymentTask{
		UserId:        "61df93c0ac601d1be8e64613",
		TransactionId: "61dfa20adebb9d4fb62b9703",
		AccountId:     "61df9b621d2c2b15a6e53ec9",
		Amount:        325,
		MetaData: &core.MetaData{
			PreferredPlanType:    core.PlanType_PLANTYPE_MIN_FEES,
			PreferredPaymentFreq: core.PaymentFrequency_PAYMENTFREQ_QUARTERLY,
		},
	}

	paymentTask, err := server.UpdatePaymentTask(ctx, &core.UpdatePaymentTaskRequest{Id: "61dfa8296c734067e6726761", PaymentTask: u})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if paymentTask.MetaData.PreferredPaymentFreq != u.MetaData.PreferredPaymentFreq {
		t.Errorf("2: Failed to fetch correct paymentTask: %+v", paymentTask)
	}
}

func TestCoreServer_DeletePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	u, err := server.ListPaymentTasks(ctx, &core.ListPaymentTaskRequest{})
	originalLen := len(u.GetPaymentTasks())

	newPaymentTask := PaymentTaskPBToDB(
		&core.PaymentTask{
			UserId:        "61df93c0ac601d1be8e64613",
			TransactionId: "61dfa20adebb9d4fb62b9703",
			AccountId:     "61df9b621d2c2b15a6e53ec9",
			Amount:        1000,
			MetaData: &core.MetaData{
				PreferredPlanType:    core.PlanType_PLANTYPE_MIN_FEES,
				PreferredPaymentFreq: core.PaymentFrequency_PAYMENTFREQ_QUARTERLY,
			},
		},
		primitive.NewObjectID(),
	)
	_, err = server.PaymentTaskDB.InsertOne(ctx, &newPaymentTask)
	if err != nil {
		t.Errorf("1: Error creating new paymentTask:: %v", err)
	}

	paymentTasks, err := server.ListPaymentTasks(ctx, &core.ListPaymentTaskRequest{})
	if err != nil {
		t.Errorf("2: An error was returned: %v", err)
	}
	newLen := len(paymentTasks.GetPaymentTasks())
	if newLen != originalLen+1 {
		t.Errorf("3: An error adding a temp paymentTask, number of paymentTasks in DB: %v", newLen)
	}

	deleted, err := server.DeletePaymentTask(ctx, &core.DeletePaymentTaskRequest{Id: newPaymentTask.ID.Hex()})
	if err != nil {
		t.Errorf("4: An error was returned: %v", err)
	}
	if deleted.Status != core.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete paymentTask: %+v\n, %+v", deleted.Status, deleted.GetPaymentTask())
	}
}

func TestCoreServer_CreateManyPaymentTasks(t *testing.T) {
	server, ctx := GenServer()
	limit := 10
	c := 0
	var res []*core.PaymentTask
	for ok := true; ok; ok = c < limit {
		task, _ := GenFakePaymentTask()
		res = append(res, task)
		c++
	}

	_, err := server.CreateManyPaymentTask(ctx, &core.CreateManyPaymentTaskRequest{PaymentTasks: res})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
}
