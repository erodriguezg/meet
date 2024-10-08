package exception

import "github.com/erodriguezg/meet/pkg/core/domain"

func NewModelAlreadyRegisteredException(model *domain.Model) error {
	return newBusinessException("model-already-registered",
		"the person was already registered a model",
		map[string]string{"personId": model.PersonId.Hex(), "nickName": model.NickName})
}

func NewModelNickNameNotAvailable(nickName string) error {
	return newBusinessException("model-nickname-not-available",
		"the model nickname is not available for register",
		map[string]string{"nickName": nickName})
}
