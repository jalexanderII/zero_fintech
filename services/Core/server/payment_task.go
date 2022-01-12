package server

import (
	"context"

	"github.com/jalexanderII/zero_fintech/services/Core/gen/core"
)

func (s CoreServer) CreatePaymentTask(ctx context.Context, in *core.CreatePaymentTaskRequest) (*core.PaymentTask, error) {
	// sername, email, password, role := in.GetUsername(), in.GetEmail(), in.GetPassword(), in.GetRole()
	// match, _ := regexp.MatchString(config.EmailRegex, email)
	// if !match {
	// 	return nil, errors.New("email validation failed")
	// }
	//
	// res, err := s.UsernameUsed(ctx, username)
	// if err != nil {
	// 	log.Printf("Error returned from UsernameUsed: %v\n", err)
	// 	return nil, err
	// }
	// if res {
	// 	return nil, errors.New("username already taken")
	// }
	// res, err = s.EmailUsed(ctx, email)
	// if err != nil {
	// 	log.Printf("Error returned from EmailUsed: %v\n", err)
	// 	return nil, err
	// }
	// if res {
	// 	return nil, errors.New("email already used")
	// }
	//
	// pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// newUser := userDB.User{ID: primitive.NewObjectID(), Email: email, Username: username, Password: string(pw), Role: userDB.Role(role)}
	//
	// _, err = s.DB.InsertOne(ctx, newUser)
	// if err != nil {
	// 	log.Printf("Error inserting new user: %v\n", err)
	// 	return nil, err
	// }
	return nil, nil
}
func (s CoreServer) GetPaymentTask(ctx context.Context, in *core.GetPaymentTaskRequest) (*core.PaymentTask, error) {
	// var user userDB.User
	// id, err := primitive.ObjectIDFromHex(in.GetId())
	// if err != nil {
	// 	return nil, err
	// }
	//
	// filter := bson.D{{"_id", id}}
	// err = s.DB.FindOne(ctx, filter).Decode(&user)
	// if err != nil {
	// 	return nil, err
	// }
	// return userDBToPB(&user), nil
	return nil, nil
}
func (s CoreServer) ListPaymentTasks(ctx context.Context, in *core.ListPaymentTaskRequest) (*core.ListPaymentTaskResponse, error) {
	// var results []userDB.User
	// cursor, err := s.DB.Find(ctx, bson.D{})
	// if err = cursor.All(ctx, &results); err != nil {
	// 	s.l.Error("[DB] Error getting all users", "error", err)
	// 	return nil, err
	// }
	// res := make([]*userPB.User, len(results))
	// for idx, user := range results {
	// 	res[idx] = userDBToPB(&user)
	// }
	// return &userPB.ListUserResponse{Users: res}, nil
	return nil, nil
}
func (s CoreServer) UpdatePaymentTask(ctx context.Context, in *core.UpdatePaymentTaskRequest) (*core.PaymentTask, error) {
	// username, email := in.GetUser().GetUsername(), in.GetUser().GetEmail()
	// id, err := primitive.ObjectIDFromHex(in.GetId())
	// if err != nil {
	// 	return nil, err
	// }
	// filter := bson.D{{"_id", id}}
	// update := bson.D{{"$set", bson.D{{"username", username}, {"email", email}}}}
	// _, err = s.DB.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	return nil, err
	// }
	// var user userDB.User
	// err = s.DB.FindOne(ctx, filter).Decode(&user)
	// return userDBToPB(&user), nil
	return nil, nil
}
func (s CoreServer) DeletePaymentTask(ctx context.Context, in *core.DeletePaymentTaskRequest) (*core.DeletePaymentTaskResponse, error) {
	// id, err := primitive.ObjectIDFromHex(in.GetId())
	// if err != nil {
	// 	return nil, err
	// }
	// filter := bson.D{{"_id", id}}
	// _, err = s.DB.DeleteOne(ctx, filter)
	// if err != nil {
	// 	return nil, err
	// }
	// var user userDB.User
	// err = s.DB.FindOne(ctx, filter).Decode(&user)
	// return &userPB.DeleteUserResponse{Status: userPB.STATUS_SUCCESS, User: userDBToPB(&user)}, nil
	return nil, nil
}
