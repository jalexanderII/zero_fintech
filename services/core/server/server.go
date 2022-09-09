package server

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jalexanderII/zero_fintech/gen/Go/common"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"github.com/jalexanderII/zero_fintech/gen/Go/notification"
	"github.com/jalexanderII/zero_fintech/gen/Go/planning"
	"github.com/jalexanderII/zero_fintech/services/auth/config/middleware"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// CoreServer is the server for the CoreService, it will connect to its own mongodb database and will be reachable via
// grpc from microservices and via grpc proxy for clients
type CoreServer struct {
	core.UnimplementedCoreServer
	// Database collections
	PaymentTaskDB mongo.Collection
	AccountDB     mongo.Collection
	TransactionDB mongo.Collection
	UserDB        mongo.Collection
	// authentication manager
	jwtm *middleware.JWTManager
	// clients
	planningClient     planning.PlanningClient
	notificationClient notification.NotificationClient
	// custom logger
	l *logrus.Logger
}

func NewCoreServer(pdb mongo.Collection, adb mongo.Collection, tdb mongo.Collection,
	udb mongo.Collection, jwtm *middleware.JWTManager, pc planning.PlanningClient, l *logrus.Logger,
) *CoreServer {
	return &CoreServer{
		PaymentTaskDB: pdb, AccountDB: adb, TransactionDB: tdb,
		UserDB: udb, jwtm: jwtm, planningClient: pc, l: l,
	}
}

func (s CoreServer) GetPaymentPlan(ctx context.Context, in *core.GetPaymentPlanRequest) (*common.PaymentPlanResponse, error) {
	// create payment task from user inputs
	paymentTasks := make([]*common.PaymentTask, len(in.GetAccountInfo()))

	for idx, item := range in.GetAccountInfo() {
		task := &common.PaymentTask{
			UserId:    in.UserId,
			AccountId: item.AccountId,
			Amount:    item.Amount,
		}
		paymentTasks[idx] = task
	}

	// save payment tasks to DB
	listOfIds, err := s.CreateManyPaymentTask(ctx, &common.CreateManyPaymentTaskRequest{PaymentTasks: paymentTasks})
	if err != nil {
		s.l.Error("[PaymentTask] Error creating PaymentTasks", "error", err)
		return nil, err
	}

	for idx, id := range listOfIds.GetPaymentTaskIds() {
		pt, _ := s.GetPaymentTask(ctx, &common.GetPaymentTaskRequest{Id: id})
		paymentTasks[idx] = pt
	}

	// send payment tasks to planning to get payment plans
	res, err := s.planningClient.CreatePaymentPlan(ctx, &planning.CreatePaymentPlanRequest{PaymentTasks: paymentTasks, MetaData: in.GetMetaData(), SavePlan: in.GetSavePlan()})
	if err != nil {
		return nil, err
	}
	return &common.PaymentPlanResponse{PaymentPlans: res.GetPaymentPlans()}, nil
}

func (s CoreServer) NotifyUsersUpcomingPaymentActions(ctx context.Context, in *notification.NotifyUsersUpcomingPaymentActionsRequest) (*notification.NotifyUsersUpcomingPaymentActionsResponse, error) {
	now := ptypes.TimestampProto(time.Now())
	upcomingPaymentActionsAllUsers, err := s.planningClient.GetAllUpcomingPaymentActions(ctx, &planning.ListPaymentPlanRequest{date: now})
	if err != nil {
		s.l.Error("[PaymentPlan] Error listing upcoming PaymentActions", "error", err)
		return nil, err
	}
	userIds := upcomingPaymentActionsAllUsers.GetUserIds()
	paymentActions := upcomingPaymentActionsAllUsers.GetPaymentActions()

	// create map of UserID -> AccID -> Liability
	userAccLiabilities := make(map[string]map[string]float64)
	for idx := range upcomingPaymentActionsAllUsers {
		_, created := userAccLiabilities[userIds[idx]]
		if !created {
			userAccLiabilities[userIds[idx]] = make(map[string]float64)
		}
		userAccLiabilities[userIds[idx]][paymentActions[idx].GetAccountId()] += paymentActions[idx].GetAmount()
	}

	// creates map of how to inform users
	userNotify := make(map[string]string)
	for userId, accLiab := range userAccLiabilities {
		totalLiab := 0.0
		for _, liab := range accLiab {
			totalLiab += liab
		}

		// TODO: call BFF instead to get amount of debit available
		userAccs, err := s.ListUserAccounts(ctx, &core.ListUserAccountsRequest{UserId: userId})
		if err != nil {
			s.l.Error("[Accounts] Error listing accounts for user", userId, "error", err)
			return nil, err
		}
		totalDebit, err := s.GetDebitAccountBalance(ctx, &core.GetDebitAccountBalanceRequest{UserId: userId})
		if err != nil {
			s.l.Error("[Accounts] Error getting debit balance for user", userId, "error", err)
			return nil, err
		}

		if totalDebit < totalLiab {
			userNotify[userId] = fmt.Sprintf("You are missing %v for tomorrows upcoming total payment of %v", totalLiab-totalDebit, totalLiab)
		} else {
			userNotify[userId] = fmt.Sprintf("You are all setup for tomorrows total payment of %v", totalLiab)
		}
		for accId, liab := range accLiab {
			accName := ""
			for acc := range userAccs.GetAccounts() {
				if acc.GetAccountId() == accId {
					accName = acc.GetName()
					break
				}
			}
			userNotify[userId] += fmt.Sprintf("\n%v: %v", accName, liab)
		}
	}

	// send notifications to the appropriate user
	for userId, message := range userNotify {
		phoneNumber := s.GetUser(ctx, &core.GetUserRequest{Id: userId}).GetPhoneNumber()
		_, err := s.notificationClient.SendSMS(ctx, &notification.SendSMSRequest{PhoneNumber: phoneNumber, Message: message})
		if err != nil {
			s.l.Error("[Notification] Failed to notify user", userId, "error", err)
			return nil, err
		}
	}
	return nil, nil
}

func (s CoreServer) GetWaterfallOverview(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.WaterfallOverviewResponse, error) {
	return s.planningClient.GetWaterfallOverview(ctx, in)
}

func (s CoreServer) GetAmountPaidPercentage(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.GetAmountPaidPercentageResponse, error) {
	return s.planningClient.GetAmountPaidPercentage(ctx, in)
}

func (s CoreServer) GetPercentageCoveredByPlans(ctx context.Context, in *planning.GetUserOverviewRequest) (*planning.GetPercentageCoveredByPlansResponse, error) {
	return s.planningClient.GetPercentageCoveredByPlans(ctx, in)
}

func (s CoreServer) ListUserPaymentPlans(ctx context.Context, in *common.ListUserPaymentPlansRequest) (*common.ListPaymentPlanResponse, error) {
	return s.planningClient.ListUserPaymentPlans(ctx, in)
}
