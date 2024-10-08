package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuyPackService interface {
	GetPaymentClientData() (map[string]any, error)

	GetPackBuyDetails(modelNickName string, packNumber int) (*dto.PackBuyDetailDto, error)

	CreateBuyPackOrder(buyerPersonId string, modelNickName string, packNumber int) (string, error)

	CapturePackPayment(orderID string) error
}

type domainBuyPackService struct {
	modelService           ModelService
	personService          PersonService
	packService            PackService
	ownedResourceService   OwnedResourceService
	paymentClient          repository.PaymentClientRepository
	paymentOrderRepository repository.PaymentOrderRepository
}

func NewDomainBuyPackService(personService PersonService,
	modelService ModelService,
	packService PackService,
	ownedResourceService OwnedResourceService,
	paymentClient repository.PaymentClientRepository,
	paymentOrderRepository repository.PaymentOrderRepository) BuyPackService {
	return &domainBuyPackService{
		modelService,
		personService,
		packService,
		ownedResourceService,
		paymentClient,
		paymentOrderRepository,
	}
}

func (port *domainBuyPackService) GetPaymentClientData() (map[string]any, error) {
	return port.paymentClient.GetClientData()
}

func (port *domainBuyPackService) GetPackBuyDetails(modelNickName string, packNumber int) (*dto.PackBuyDetailDto, error) {
	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, fmt.Errorf("model not found")
	}
	pack, err := port.packService.FindActivePackByModelIdAndPackNumber(model.Id.Hex(), packNumber)
	if err != nil {
		return nil, err
	}
	if pack.DollarValue == nil {
		return nil, fmt.Errorf("the pack does not have dollar value")
	}

	dto := dto.PackBuyDetailDto{
		ModelNickName:   modelNickName,
		PackTitle:       pack.Title,
		PackDollarValue: *pack.DollarValue,
	}
	return &dto, nil
}

func (port *domainBuyPackService) CreateBuyPackOrder(personId string, modelNickName string, packNumber int) (string, error) {

	person, err := port.personService.FindById(personId)
	if err != nil {
		return "", err
	}
	if person == nil {
		return "", fmt.Errorf("person id %s not found for CreateBuyPackOrder", personId)
	}

	model, err := port.modelService.FindModelByNickName(modelNickName)
	if err != nil {
		return "", err
	}
	if model == nil {
		return "", fmt.Errorf("model %s not found for CreateBuyPackOrder", modelNickName)
	}

	pack, err := port.packService.FindActivePackByModelIdAndPackNumber(model.Id.Hex(), packNumber)
	if err != nil {
		return "", err
	}
	if pack == nil {
		return "", fmt.Errorf("pack not found for CreateBuyPackOrder")
	}
	if !pack.Published {
		return "", fmt.Errorf("pack is not published yet")
	}

	personHasPack, err := port.ownedResourceService.PersonHasPack(personId, pack.Id.Hex())
	if err != nil {
		return "", err
	}
	if personHasPack {
		return "", fmt.Errorf("the person id %s already has pack id %s", personId, pack.Id.Hex())
	}

	packDollarValue := pack.DollarValue
	if packDollarValue == nil {
		return "", fmt.Errorf("the pack does not have value")
	}

	orderId, err := port.paymentClient.CreateOrder(*packDollarValue, "USD")
	if err != nil {
		return "", err
	}

	personObjectId, err := primitive.ObjectIDFromHex(personId)
	if err != nil {
		return "", err
	}
	packObjectId, err := primitive.ObjectIDFromHex(pack.Id.Hex())
	if err != nil {
		return "", err
	}

	paymentOrder := domain.PaymentOrder{
		OrderId:            orderId,
		PersonId:           personObjectId,
		PackId:             packObjectId,
		ModelId:            pack.ModelId,
		PaymentDollarValue: *pack.DollarValue,
		CreatedAt:          time.Now(),
	}

	paymentOrderSaved, err := port.paymentOrderRepository.SavePaymentOrder(&paymentOrder)
	if err != nil {
		return "", err
	}

	return paymentOrderSaved.OrderId, nil
}

func (port *domainBuyPackService) CapturePackPayment(orderID string) error {

	paymentOrder, err := port.paymentOrderRepository.FindByOrderId(orderID)
	if err != nil {
		return err
	}
	if paymentOrder == nil {
		return fmt.Errorf("no payment order found for id: %s", orderID)
	}
	if paymentOrder.CapturedAt != nil {
		return fmt.Errorf("the payment order id: %s was already capture", orderID)
	}

	paymentDetails, err := port.paymentClient.CapturePayment(orderID)
	if err != nil {
		return err
	}

	paymentOrder.PaymentDetails = paymentDetails
	presentTime := time.Now()
	paymentOrder.CapturedAt = &presentTime

	_, err = port.paymentOrderRepository.SavePaymentOrder(paymentOrder)
	if err != nil {
		return err
	}

	err = port.ownedResourceService.AddPackToPerson(paymentOrder.PersonId.Hex(), paymentOrder.PackId.Hex())
	if err != nil {
		return err
	}

	return nil
}
