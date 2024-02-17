package test

import (
	"context"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/types"
	"github.com/stretchr/testify/mock"
)

// Mock indexer service
type MockIndexerService struct {
	mock.Mock
}

// IndexERC20Transfer mocks the IndexERC20Transfer method
func (m *MockIndexerService) IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress) error {
	return nil
}

// IndexOrderCreated mocks the IndexOrderCreated method
func (m *MockIndexerService) IndexOrderCreated(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

// IndexOrderSettled mocks the IndexOrderSettled method
func (m *MockIndexerService) IndexOrderSettled(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

// IndexOrderRefunded mocks the IndexOrderRefunded method
func (m *MockIndexerService) IndexOrderRefunded(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

// HandleReceiveAddressValidity mocks the HandleReceiveAddressValidity method
func (m *MockIndexerService) HandleReceiveAddressValidity(ctx context.Context, receiveAddress *ent.ReceiveAddress, paymentOrder *ent.PaymentOrder) error {
	return nil
}

// Mock order service
type MockOrderService struct {
	mock.Mock
}

// CreateOrder mocks the CreateOrder method
func (m *MockOrderService) CreateOrder(ctx context.Context, orderID uuid.UUID) error {
	return nil
}

// RefundOrder mocks the RefundOrder method
func (m *MockOrderService) RefundOrder(ctx context.Context, orderID string) error {
	return nil
}

// RevertOrder mocks the RevertOrder method
func (m *MockOrderService) RevertOrder(ctx context.Context, order *ent.PaymentOrder) error {
	return nil
}

// SettleOrder mocks the SettleOrder method
func (m *MockOrderService) SettleOrder(ctx context.Context, orderID uuid.UUID) error {
	return nil
}

// GetSupportedInstitutions mocks the GetSupportedInstitutions method
func (m *MockOrderService) GetSupportedInstitutions(ctx context.Context, client types.RPCClient, currencyCode string) ([]types.Institution, error) {
	return nil, nil
}
