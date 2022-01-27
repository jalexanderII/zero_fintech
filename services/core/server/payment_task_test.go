package server

import (
	"testing"

	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreServer_CreatePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	equinoxPT := &common.PaymentTask{
		UserId:    "61df93c0ac601d1be8e64613",
		AccountId: "61df9b621d2c2b15a6e53ec9",
		Amount:    325,
	}

	paymentTask, err := server.CreatePaymentTask(ctx, &common.CreatePaymentTaskRequest{PaymentTask: equinoxPT})
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

	paymentTask, err := server.GetPaymentTask(ctx, &common.GetPaymentTaskRequest{Id: "61dfa8296c734067e6726761"})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if paymentTask.Amount != 325 {
		t.Errorf("2: Failed to fetch correct paymentTask: %+v", paymentTask)
	}
}

func TestCoreServer_ListPaymentTasks(t *testing.T) {
	server, ctx := GenServer()

	paymentTasks, err := server.ListPaymentTasks(ctx, &common.ListPaymentTaskRequest{})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if len(paymentTasks.PaymentTasks) < 1 {
		t.Errorf("2: Failed to fetch realtors: %+v", paymentTasks.PaymentTasks[0])
	}
}

func TestCoreServer_UpdatePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	u := &common.PaymentTask{
		UserId:    "61df93c0ac601d1be8e64613",
		AccountId: "61df9b621d2c2b15a6e53ec9",
		Amount:    325,
	}

	paymentTask, err := server.UpdatePaymentTask(ctx, &common.UpdatePaymentTaskRequest{Id: "61dfa8296c734067e6726761", PaymentTask: u})
	if err != nil {
		t.Errorf("1: An error was returned: %v", err)
	}
	if paymentTask.Amount != u.Amount {
		t.Errorf("2: Failed to fetch correct paymentTask: %+v", paymentTask)
	}
}

func TestCoreServer_DeletePaymentTask(t *testing.T) {
	server, ctx := GenServer()

	u, _ := server.ListPaymentTasks(ctx, &common.ListPaymentTaskRequest{})
	originalLen := len(u.GetPaymentTasks())

	newPaymentTask := PaymentTaskPBToDB(
		&common.PaymentTask{
			UserId:    "61df93c0ac601d1be8e64613",
			AccountId: "61df9b621d2c2b15a6e53ec9",
			Amount:    1000,
		},
		primitive.NewObjectID(),
	)
	_, err := server.PaymentTaskDB.InsertOne(ctx, &newPaymentTask)
	if err != nil {
		t.Errorf("1: Error creating new paymentTask:: %v", err)
	}

	paymentTasks, err := server.ListPaymentTasks(ctx, &common.ListPaymentTaskRequest{})
	if err != nil {
		t.Errorf("2: An error was returned: %v", err)
	}
	newLen := len(paymentTasks.GetPaymentTasks())
	if newLen != originalLen+1 {
		t.Errorf("3: An error adding a temp paymentTask, number of paymentTasks in DB: %v", newLen)
	}

	deleted, err := server.DeletePaymentTask(ctx, &common.DeletePaymentTaskRequest{Id: newPaymentTask.ID.Hex()})
	if err != nil {
		t.Errorf("4: An error was returned: %v", err)
	}
	if deleted.Status != common.DELETE_STATUS_DELETE_STATUS_SUCCESS {
		t.Errorf("5: Failed to delete paymentTask: %+v\n, %+v", deleted.Status, deleted.GetPaymentTask())
	}
}

func TestCoreServer_CreateManyPaymentTasks(t *testing.T) {
	server, ctx := GenServer()
	limit := 10
	c := 0
	var res []*common.PaymentTask
	for ok := true; ok; ok = c < limit {
		task, _ := GenFakePaymentTask()
		res = append(res, task)
		c++
	}

	_, err := server.CreateManyPaymentTask(ctx, &common.CreateManyPaymentTaskRequest{PaymentTasks: res})
	if err != nil {
		t.Errorf("1: Error creating new paymentTask: %v", err)
	}
}
