package transaction

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/sirupsen/logrus"
)

type TransactionRepository interface {
	CreateCharge(ctx context.Context, orderID string, grossAmount int64, paymentType string, params map[string]interface{}) (string, string, error)
	GetStatus(ctx context.Context, transactionID string) (string, error)
	VerifyNotification(orderID, statusCode, grossAmount, signatureKey string) bool
}

type MidtransRepository struct {
	core      *coreapi.Client
	serverKey string
}

func NewMidtransRepository() (*MidtransRepository, error) {
	key := os.Getenv("MIDTRANS_SERVER_KEY")
	env := os.Getenv("MIDTRANS_ENV")

	if key == "" {
		return nil, fmt.Errorf("midtrans not configured: missing SERVER_KEY")
	}

	c := coreapi.Client{}
	if env == "production" {
		c.New(key, midtrans.Production)
	} else {
		c.New(key, midtrans.Sandbox)
	}

	return &MidtransRepository{
		core:      &c,
		serverKey: key,
	}, nil
}

func (m *MidtransRepository) CreateCharge(ctx context.Context, orderID string, grossAmount int64, paymentType string, params map[string]interface{}) (string, string, error) {
	_ = ctx // ctx unused - Midtrans SDK doesn't support context
	// Log unknown payment types but allow them
	knownTypes := map[string]bool{
		"credit_card": true, "bank_transfer": true, "echannel": true,
		"gopay": true, "shopeepay": true, "qris": true,
	}
	if !knownTypes[paymentType] {
		logrus.WithField("payment_type", paymentType).Warn("Unknown payment type, proceeding anyway")
	}

	req := &coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType(paymentType),
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
	}

	resp, err := m.core.ChargeTransaction(req)
	if err != nil {
		logrus.WithError(err).WithField("order_id", orderID).Error("Midtrans charge failed")
		return "", "", fmt.Errorf("payment gateway error: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"order_id":           orderID,
		"transaction_id":     resp.TransactionID,
		"transaction_status": resp.TransactionStatus,
		"fraud_status":       resp.FraudStatus,
	}).Info("Midtrans charge created")

	var redirect string
	if resp.RedirectURL != "" {
		redirect = resp.RedirectURL
	}

	return resp.TransactionID, redirect, nil
}

func (m *MidtransRepository) GetStatus(ctx context.Context, transactionID string) (string, error) {
	_ = ctx // ctx unused - Midtrans SDK doesn't support context
	resp, err := m.core.CheckTransaction(transactionID)
	if err != nil {
		logrus.WithError(err).WithField("transaction_id", transactionID).Error("Midtrans status check failed")
		return "", fmt.Errorf("failed to check payment status: %w", err)
	}
	return resp.TransactionStatus, nil
}

func (m *MidtransRepository) VerifyNotification(orderID, statusCode, grossAmount, signatureKey string) bool {
	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + m.serverKey))
	expected := hex.EncodeToString(sum[:])
	return strings.EqualFold(expected, signatureKey)
}

type MockTransactionRepository struct {
	Charges []MockCharge
}

type MockCharge struct {
	OrderID     string
	Amount      int64
	PaymentType string
	Params      map[string]interface{}
}

func (m *MockTransactionRepository) CreateCharge(ctx context.Context, orderID string, grossAmount int64, paymentType string, params map[string]interface{}) (string, string, error) {
	_ = ctx // ctx unused in mock
	m.Charges = append(m.Charges, MockCharge{
		OrderID:     orderID,
		Amount:      grossAmount,
		PaymentType: paymentType,
		Params:      params,
	})
	return "mock-tx-" + orderID, "https://mock-payment.com/redirect", nil
}

func (m *MockTransactionRepository) GetStatus(ctx context.Context, transactionID string) (string, error) {
	_ = ctx // ctx unused in mock
	return "paid", nil // Always paid for testing
}

func (m *MockTransactionRepository) VerifyNotification(orderID, statusCode, grossAmount, signatureKey string) bool {
	return true // Always valid for testing
}

// MapStatusToInternal maps Midtrans status to internal status
func MapStatusToInternal(midtransStatus string) string {
	switch midtransStatus {
	case "capture", "settlement":
		return "paid"
	case "pending":
		return "pending"
	case "deny", "cancel", "expire", "failure":
		return "failed"
	default:
		return midtransStatus
	}
}