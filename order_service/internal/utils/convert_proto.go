package utils

import (
	"order-service-system/order_service/internal/models"
	orderpb "order-service-system/proto/order"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertToProto(doc models.Order) *orderpb.Order {
	items := make([]*orderpb.OrderItem, 0, len(doc.Items))
	for _, item := range doc.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	status, ok := orderpb.OrderStatus_value[doc.Status]
	if !ok {
		status = int32(orderpb.OrderStatus_ORDER_STATUS_UNSPECIFIED)
	}

	return &orderpb.Order{
		OrderId:     doc.OrderID,
		UserId:      doc.UserID,
		Items:       items,
		TotalAmount: doc.TotalAmount,
		Status:      orderpb.OrderStatus(status),
		CreatedAt:   timestamppb.New(doc.CreatedAt),
	}
}
