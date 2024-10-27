package enumdef

import (
	"fmt"

	"github.com/todennus/proto/gen/service/dto"
)

type ConfidentialRequirementType int

const (
	CRTRequire ConfidentialRequirementType = iota
	CRTNotRequire
	CRTDependOnType
)

func ConfidentialRequirementTypeFromGRPC(req dto.OAuth2ClientConfidentialRequirement) ConfidentialRequirementType {
	switch req {
	case dto.OAuth2ClientConfidentialRequirement_DEPEND:
		return CRTDependOnType
	case dto.OAuth2ClientConfidentialRequirement_NOT_REQUIRE:
		return CRTNotRequire
	case dto.OAuth2ClientConfidentialRequirement_REQUIRE:
		return CRTRequire
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}

func ConfidentialRequirementTypeToGRPC(req ConfidentialRequirementType) dto.OAuth2ClientConfidentialRequirement {
	switch req {
	case CRTDependOnType:
		return dto.OAuth2ClientConfidentialRequirement_DEPEND
	case CRTRequire:
		return dto.OAuth2ClientConfidentialRequirement_REQUIRE
	case CRTNotRequire:
		return dto.OAuth2ClientConfidentialRequirement_NOT_REQUIRE
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}
