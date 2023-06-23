package user_test

import (
	"errors"
	"testing"

	mock_notification "go-clean/src/business/domain/mock/notification"
	mock_user "go-clean/src/business/domain/mock/user"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/user"
	mock_auth "go-clean/src/lib/tests/mock/auth"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	notificationMock := mock_notification.NewMockInterface(ctrl)
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	mockParams := entity.CreateUserParam{
		Email:    "r@g.com",
		Password: "password",
	}

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "r@g.com",
		Password: string(hashPass),
	}

	mockNotificationWelcome := entity.Notification{
		UserID:      1,
		Description: "Selamat Datang di <b>Office Buddy</b>",
		Status:      entity.WelcomeStatus,
		IsRead:      false,
	}

	mockNotificationFirstOrder := entity.Notification{
		UserID:      1,
		Description: "Mulai lakukan <b>pemesanan kantor atau co-working space</b> pertama kamu!",
		Status:      entity.FirstOrderStatus,
		IsRead:      false,
	}

	u := user.Init(userMock, nil, notificationMock)

	type mockfields struct {
		user         *mock_user.MockInterface
		notification *mock_notification.MockInterface
	}

	mocks := mockfields{
		user:         userMock,
		notification: notificationMock,
	}

	type args struct {
		params entity.CreateUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to create user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Email: "r@g.com",
			},
			wantErr: true,
		},
		{
			name: "failed to create notification",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, nil)
				mock.notification.EXPECT().Create(gomock.Any()).Return(mockNotificationWelcome, assert.AnError)
				mock.notification.EXPECT().Create(gomock.Any()).Return(mockNotificationFirstOrder, assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Email: "r@g.com",
			},
			wantErr: false,
		},
		{
			name: "all ok",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, nil)
				mock.notification.EXPECT().Create(gomock.Any()).Return(mockNotificationWelcome, nil)
				mock.notification.EXPECT().Create(gomock.Any()).Return(mockNotificationFirstOrder, nil)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Email: "r@g.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Create(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Email, got.Email)
		})
	}
}

func Test_user_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	authMock := mock_auth.NewMockInterface(ctrl)

	mockParams := entity.LoginUserParam{
		Email:    "r@g.com",
		Password: "password",
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Email:    "r@g.com",
		Password: string(hashPass),
	}

	mockToken := "mockToken"

	u := user.Init(userMock, authMock, nil)

	type mockfields struct {
		user *mock_user.MockInterface
		auth *mock_auth.MockInterface
	}

	mocks := mockfields{
		user: userMock,
		auth: authMock,
	}

	type args struct {
		params entity.LoginUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "failed to find user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().GetByEmail(arg.params.Email).Return(entity.User{}, errors.New("user not found"))
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().GetByEmail(arg.params.Email).Return(entity.User{}, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "password incorrect",
			mockFunc: func(mock mockfields, arg args) {
				mockUserResultWithWrongPassword := mockUserResult
				mockUserResultWithWrongPassword.Password = "wrongPassword"
				mock.user.EXPECT().GetByEmail(arg.params.Email).Return(mockUserResultWithWrongPassword, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed to generate token",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().GetByEmail(arg.params.Email).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("failed to generate token"))
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().GetByEmail(arg.params.Email).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return(mockToken, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    mockToken,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Login(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
