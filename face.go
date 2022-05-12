package superaipro

import "context"

type IClient interface {
	Send(context.Context, Request) (Result, error)
	GetResult(context.Context, string) (Result, error)
	GetUser() (User, error)
	GetWallet() (Wallet, error)
	Solve(context.Context, Request) (Result, error)
	Identify(Request) (IdentifyResult, error)
}
