package payment

import (
	"bwastartup/user"
	"strconv"

	//tidak bisa pakai dash, jadi dibuat alias
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	//ini karena dari midtrans, membutuhkan data transaksi dan user
	//awalnya GetToken, seperti dokumentasi midtrans
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-0uMntmE_ygUraV2LSQH4xWpH"
	midclient.ClientKey = "SB-Mid-client-syR5mcKQ8Q01RxVD"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
