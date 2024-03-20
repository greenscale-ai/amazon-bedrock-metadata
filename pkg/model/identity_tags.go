package model

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/iam"
	"strings"
)

type IdentityTagsBuilder struct {
	iamClient       *iam.IAM
	entityTagsCache map[string][]*iam.Tag
}

func NewIdentityTagsBuilder(iamClient *iam.IAM) *IdentityTagsBuilder {
	return &IdentityTagsBuilder{
		iamClient:       iamClient,
		entityTagsCache: make(map[string][]*iam.Tag),
	}
}

func (i *IdentityTagsBuilder) parseIamEntity(identityArn string) (string, string) {
	parsedARN, err := arn.Parse(identityArn)
	if err != nil {
		return "Unknown", "Unknown"
	}

	entityParts := strings.SplitN(parsedARN.Resource, "/", 2)
	if len(entityParts) == 2 {
		return entityParts[0], entityParts[1] // Return both the IAM entity type and name
	}

	return "Unknown", "Unknown"
}

func (i *IdentityTagsBuilder) GetIdentityTags(identity string) (tags []*iam.Tag, err error) {

	iamEntityType, entityName := i.parseIamEntity(identity)

	val, ok := i.entityTagsCache[fmt.Sprintf("%s:%s", iamEntityType, entityName)]

	if ok {
		return val, nil
	}

	switch iamEntityType {
	case "user":
		result, err := i.iamClient.ListUserTags(&iam.ListUserTagsInput{
			UserName: aws.String(entityName),
		})
		if err != nil {
			return nil, err
		}
		tags = result.Tags
		break

	case "role":
		result, err := i.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: aws.String(entityName),
		})
		if err != nil {
			return nil, err
		}
		tags = result.Tags
		break

	case "assumed-role":
		roleName := strings.Split(entityName, "/")[0]

		result, err := i.iamClient.ListRoleTags(&iam.ListRoleTagsInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return nil, err
		}
		tags = result.Tags
		break

	case "instance-profile":
		result, err := i.iamClient.ListInstanceProfileTags(&iam.ListInstanceProfileTagsInput{
			InstanceProfileName: aws.String(entityName),
		})
		if err != nil {
			return nil, err
		}
		tags = result.Tags
		break

	case "saml-provider":
		result, err := i.iamClient.ListSAMLProviderTags(&iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(identity),
		})
		if err != nil {
			return nil, err
		}
		tags = result.Tags
		break

	default:
		return nil, errors.New("unsupported IAM entity type")
	}

	i.entityTagsCache[fmt.Sprintf("%s:%s", iamEntityType, entityName)] = tags
	return tags, nil
}
