package enumdef

import (
	"fmt"

	"github.com/todennus/proto/gen/service/dto"
)

type OAuth2ClientConfidentialRequirement int

const (
	CRTRequire OAuth2ClientConfidentialRequirement = iota
	CRTNotRequire
	CRTDependOnType
)

func OAuth2ClientConfidentialRequirementTypeFromGRPC(req dto.OAuth2ClientConfidentialRequirement) OAuth2ClientConfidentialRequirement {
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

func OAuth2ClientConfidentialRequirementTypeToGRPC(req OAuth2ClientConfidentialRequirement) dto.OAuth2ClientConfidentialRequirement {
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
