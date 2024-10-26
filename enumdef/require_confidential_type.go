package enumdef

import (
	"fmt"

	"github.com/todennus/proto/gen/service/dto"
)

type ConfidentialRequirementType int

const (
	RequireConfidential ConfidentialRequirementType = iota
	NotRequireConfidential
	DependOnClientConfidential
)

func ConfidentialRequirementTypeFromGRPC(req dto.OAuth2ClientConfidentialRequirement) ConfidentialRequirementType {
	switch req {
	case dto.OAuth2ClientConfidentialRequirement_DEPEND:
		return DependOnClientConfidential
	case dto.OAuth2ClientConfidentialRequirement_NOT_REQUIRE:
		return NotRequireConfidential
	case dto.OAuth2ClientConfidentialRequirement_REQUIRE:
		return RequireConfidential
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}

func ConfidentialRequirementTypeToGRPC(req ConfidentialRequirementType) dto.OAuth2ClientConfidentialRequirement {
	switch req {
	case DependOnClientConfidential:
		return dto.OAuth2ClientConfidentialRequirement_DEPEND
	case RequireConfidential:
		return dto.OAuth2ClientConfidentialRequirement_REQUIRE
	case NotRequireConfidential:
		return dto.OAuth2ClientConfidentialRequirement_NOT_REQUIRE
	default:
		panic(fmt.Sprintf("invalid requirement %d", req))
	}
}
