from typing import List
import sys, os
from bson.objectid import ObjectId
from google.protobuf.timestamp_pb2 import Timestamp
from Python.common.common_pb2 import PLAN_TYPE_OPTIM_CREDIT_SCORE
# make gen/Python importable by import Python.X
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'gen')))
from Python.common.payment_task_pb2 import PaymentTask
from Python.common.common_pb2 import PAYMENT_FREQUENCY_MONTHLY
from Python.planning.payment_plan_pb2 import PAYMENT_ACTION_STATUS_PENDING, PAYMENT_STATUS_CURRENT, PaymentAction, PaymentPlan

class PaymentPlanBuilder:

    def createPaymentPlan(self, paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """ Creates a PaymentPlan protobuf from given list of PaymentTasks protos."""
        return self._mock_payment_plan(paymentTasks)

    def _mock_payment_plan(self, paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """A mock payment plan for end to end testing"""
        ids: List[str] = []
        plans: List[PaymentPlan] = []
        actions: List[PaymentAction] = []
        total = 0
        pb_timestamp = Timestamp()
        for task in paymentTasks:
            ids.append(task.payment_task_id)
            total += task.amount
            a = PaymentAction(
                account_id=task.account_id,
                amount=task.amount,
                transaction_date=pb_timestamp.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            )
            actions.append(a)
        bson_id = ObjectId()
        plan = PaymentPlan(
            payment_plan_id=str(bson_id),
            user_id=paymentTasks[0].user_id,
            payment_task_id=ids, 
            timeline=12,
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            amount_per_payment=total/12,
            plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
            end_date=pb_timestamp.GetCurrentTime(),
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=actions,
        )
        plans.append(plan)
        return plans   
            

# def CreatePaymentPlan(self, request, context: grpc.ServicerContext):
#         """ Creates PaymentPlan for given request.

#         Doesn't incorporate current PaymentPlans and assuems that one always
#         has to pay the minimum on all accounts in this plan.

#         Args:
#             request ([type]): [description]
#             context (grpc.ServicerContext): [description]

#         Returns:
#             [type]: [description]
#         """
#         user_id = None
#         payment_task_ids = []
#         transaction_ids = []
#         account_ids = []
#         pref_payment_freqs = []
#         pref_plan_types = []
#         pref_timelines = []
#         for _payment_task in request.payment_tasks:
#             if user_id is None:
#                 user_id = _payment_task.user_id
#             payment_task_ids.append(_payment_task.payment_task_id)
#             transaction_ids.append(_payment_task.transaction_id)
#             account_ids.append(_payment_task.account_id)
#             pref_payment_freqs.append(_payment_task.preferred_payment_freq)
#             pref_plan_types.append(_payment_task.preferred_plan_type)
#             pref_timelines.append(_payment_task.preferred_timeline)
#         payment_task_ids = np.array(payment_task_ids)
#         transaction_ids = np.array(transaction_ids)
#         account_ids = np.array(account_ids)

#         payment_plans = []
#         for _pref_payment_freq, _pref_payment_type, _pref_timeline in set(zip(pref_payment_freqs, pref_plan_types, pref_timelines)):
#             mask = pref_payment_freqs == _pref_payment_freq
#             payment_plans.append(self._createPaymentPlan(payment_task_ids=payment_task_ids[mask], transaction_ids=transaction_ids[mask], user_id=user_id,
#                 account_ids=account_ids[mask], pref_payment_freq=_pref_payment_freq, pref_plan_type=_pref_payment_type, pref_timeline=_pref_timeline))

#         return CreatePaymentPlanResponse(payment_plans=payment_plans)
    
#     def _createPaymentPlan(self, payment_task_ids, transaction_ids, account_ids, user_id, pref_payment_freq, pref_plan_type, pref_timeline) -> PaymentPlanPB:
#         pass
        # timestamp = Timestamp()
        # timestamp.GetCurrentTime()
        # start_date = timestamp.ToDatetime()

        # channel = grpc.insecure_channel('localhost:50051')
        # stub = core_pb2_grpc.CoreStub(channel)
        # accounts = []
        # account2apr, acc2min_payment = {}, {}
        # account2credit_limits, account2current_usage_amount, account2current_usage_prcnt = {}, {}, {}
        # for _account_id in account_ids:
        #     _account = stub.GetAccount(accounts_pb2.GetAccountRequest(id=_account_id))
        #     accounts.append(_account)
        #     _apr = (_account.annual_percentage_rate.high_end + _account.annual_percentage_rate.low_end) / 2
        #     account2apr[_account_id] = _apr
        #     account2credit_limits[_account_id] = _account.credit_limit
        #     account2current_usage_amount[_account_id] = _account.current_balance + _account.pending_transactions
        #     account2current_usage_prcnt[_account_id] = (_account.current_balance + _account.pending_transactions) / _account.credit_limit
        #     acc2min_payment[_account_id] = _account.minimum_payment_due
        # # account_ids = np.array(account_ids)
        # # account2apr = np.array(account2apr)
        # # account_amounts = []
        # account2amount = defaultdict(list)
        # for _transaction_id in transaction_ids:
        #     _transaction = stub.GetTransaction(transactions_pb2.GetTransactionRequest(id=_transaction_id))
        #     _amount = _transaction.amount
        #     # account_amounts.append(_amount)
        #     account2amount[_transaction.account_id].append(_amount)
        # account2amount = {_acc_id: sum(_amounts) for _acc_id, _amounts in account2amount.items()}
        # # account_amounts = np.array(account_amounts)
        # total_amount = sum(account2amount.values())

        # pref_payment_freq_days = None
        # if pref_payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY:
        #     pref_payment_freq_days = timedelta(days=7)
        # elif pref_payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_BIWEEKLY:
        #     pref_payment_freq_days = timedelta(days=14)
        # elif pref_payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY:
        #     pref_payment_freq_days = timedelta(days=30)
        # elif pref_payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_QUARTERLY:
        #     pref_payment_freq_days = timedelta(days=90)
        # amount_per_payment = total_amount / (timedelta(days=30) * pref_timeline / pref_payment_freq_days)

        # if pref_plan_type == payment_task_pb2.PlanType.PLAN_TYPE_MIN_FEES or pref_plan_type == payment_task_pb2.PlanType.PLANTYPE_UNKNOWN:
        #     # need to sort account IDs, APR, amount by APR
        #     account_ids = sorted(account_ids)
        #     account2apr = OrderedDict(sorted(account2apr.items()))
        #     account2amount = OrderedDict(sorted(account2amount.items()))
        #     df = pd.DataFrame({'account_id': account_ids, 'account_apr': account2apr.values(), 'account_amount': account2amount.values()}).sort_values('account_apr', ascending=False)
        #     account_ids = df['account_id'].values.tolist()
        #     account_amounts = df['account_amount'].values.tolist()

        #     payment_actions = []
        #     date = shift_date_by_payment_frequency(date=start_date, payment_freq=pref_payment_freq)
        #     timestamp.FromDatetime(date)
        #     pay_this_date = 0

        #     # pay minimum payment for every account at date
        #     for i in range(len(account_ids)):
        #         _account_id = account_ids[i]
        #         _acc_amount = account_amounts[i]
        #         _amount_this_date = min(acc2min_payment[_account_id], _acc_amount)
        #         payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #         pay_this_date += _amount_this_date
        #         account_amounts[i] -= _amount_this_date
            
        #     # pay accounts off and every minimum
        #     while len(account_amounts) > 0:
        #         _amount = account_amounts.pop(0)
        #         _account_id = account_ids.pop(0)
        #         _amount_this_date = min(amount_per_payment-pay_this_date, _amount)
        #         _amount_next_dates = _amount - _amount_this_date
        #         if _amount_this_date > 0:
        #             payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #             pay_this_date += _amount_this_date
        #         if _amount_next_dates > 0:
        #             account_amounts.insert(0, _amount_next_dates)
        #             account_ids.insert(0, _account_id)
        #             pay_this_date = 0
        #             date = shift_date_by_payment_frequency(date=date, payment_freq=pref_payment_freq)
        #             timestamp.FromDatetime(date)

        #             # pay minimum payment for every account at date
        #             for i in range(len(account_ids)):
        #                 _account_id = account_ids[i]
        #                 _acc_amount = account_amounts[i]
        #                 _amount_this_date = min(acc2min_payment[_account_id], _acc_amount)
        #                 payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #                 pay_this_date += _amount_this_date
        #                 account_amounts[i] -= _amount_this_date  
        # elif pref_plan_type == payment_task_pb2.PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE:
        #     # sort dictionaries by account_id
        #     account_ids = sorted(account_ids)
        #     account2current_usage_amount = OrderedDict(sorted(account2current_usage_amount.items()))
        #     account2credit_limits = OrderedDict(sorted(account2credit_limits.items()))
        #     account2current_usage_prcnt = OrderedDict(sorted(account2current_usage_prcnt.items()))
        #     account2amount = OrderedDict(sorted(account2amount.items()))
        #     df = pd.DataFrame({'account_id': account_ids, 'current_usage_amount': account2current_usage_amount.values(), 'current_usage_prcnt': account2current_usage_prcnt.values(), 'limit': account2credit_limits.values(), 'amount': account2amount.values()}).sort_values('current_usage_prcnt', ascending=False)
        #     # df.current_usage_prcnt = df.current_usage_prcnt.round(1)
            
        #     payment_actions = []
        #     date = shift_date_by_payment_frequency(date=start_date, payment_freq=pref_payment_freq)
        #     timestamp.FromDatetime(date)
        #     pay_this_date = 0

        #     # pay minimum payment for every account at date
        #     for i in range(len(df)):
        #         row = df[i]
        #         _account_id = row['account_id']
        #         _acc_amount = row['amount']
        #         _current_usage_amount = row['current_usage_amount']
        #         _amount_this_date = min(acc2min_payment[_account_id], _acc_amount)
        #         payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #         pay_this_date += _amount_this_date
        #         df.loc[df['account_id']==_account_id, 'amount'] = _acc_amount - _amount_this_date
        #         df.loc[df['account_id']==_account_id, 'current_usage_amount'] = _current_usage_amount - _amount_this_date
        #     # update credit card utilization
        #     df['current_usage_prcnt'] = df['current_usage_amount'] / df['limit']
        #     # pay of currently highest used credit cards and then 
        #     while len(df) > 0:
        #         # pay highest paying credit cards off first until we can't pay anything more off
        #         for i in range(len(df)):
        #             if pay_this_date >= amount_per_payment:
        #                 break
        #             row = df[i]
        #             _account_id = row['account_id']
        #             _acc_amount = row['amount']
        #             _current_usage_amount = row['current_usage_amount']
        #             _amount_this_date = min(amount_per_payment-pay_this_date, _acc_amount)
        #             payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #             pay_this_date += _amount_this_date
        #             df.loc[df['account_id']==_account_id, 'amount'] = _acc_amount - _amount_this_date
        #             df.loc[df['account_id']==_account_id, 'current_usage_amount'] = _current_usage_amount - _amount_this_date
        #         # drop any accounts whose amount is 0
        #         df = df.loc[df.amount > 0]
        #         # update credit card usage
        #         df['current_usage_prcnt'] = df['current_usage_amount'] / df['amount']
        #         # move to next date
        #         pay_this_date = 0
        #         date = shift_date_by_payment_frequency(date=date, payment_freq=pref_payment_freq)
        #         timestamp.FromDatetime(date)
        #         # pay minimum payment for every account at date
        #         for i in range(len(df)):
        #             row = df[i]
        #             _account_id = row['account_id']
        #             _acc_amount = row['amount']
        #             _current_usage_amount = row['current_usage_amount']
        #             _amount_this_date = min(acc2min_payment[_account_id], _acc_amount)
        #             payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount_this_date, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
        #             pay_this_date += _amount_this_date
        #             df.loc[df['account_id']==_account_id, 'amount'] = _acc_amount - _amount_this_date
        #             df.loc[df['account_id']==_account_id, 'current_usage_amount'] = _current_usage_amount - _amount_this_date
        #         # update credit card utilization
        #         df['current_usage_prcnt'] = df['current_usage_amount'] / df['limit']
        #         # drop any accounts whose amount is 0
        #         df = df.loc[df.amount > 0]

        # return payment_plan_pb2(payment_plan_id=1e-9, user_id=user_id, payment_task_id=payment_task_ids, amount_per_payment=amount_per_payment, plan_type=payment_task_pb2.PlanType.PLAN_TYPE_MIN_FEES, end_date=timestamp, active=True, status=payment_plan_pb2.PaymentStatus.PAYMENT_STATUS_CURRENT, payment_action=payment_actions)