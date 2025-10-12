package typeid

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
)

type TypeId string

func (ti TypeId) String() string {
	return string(ti)
}

func (ti TypeId) Type() string {
	parts := strings.SplitN(string(ti), "_", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

func (ti TypeId) Id() string {
	parts := strings.SplitN(string(ti), "_", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (ti TypeId) Length() int {
	return len(string(ti))
}

func Make(typeName string) (TypeId, error) {
	if err := validateTypeName(typeName); err != nil {
		return "", err
	}

	return TypeId(fmt.Sprintf("%s_%s", typeName, strings.ToLower(ulid.Make().String()))), nil
}

func Parse(typeName string, typeId string) (TypeId, error) {
	if err := validateTypeName(typeName); err != nil {
		return "", err
	}

	typeIdParts := strings.SplitN(typeId, "_", 2)
	if len(typeIdParts) != 2 {
		return "", fmt.Errorf("type id must contain exactly one underscore separating type and ulid")
	}

	if typeIdParts[0] != typeName {
		return "", fmt.Errorf("type id prefix '%s' does not match expected type name '%s'", typeIdParts[0], typeName)
	}

	if _, err := ulid.Parse(strings.ToUpper(typeIdParts[1])); err != nil {
		return "", fmt.Errorf("invalid ulid '%s': %w", typeIdParts[1], err)
	}

	return TypeId(typeId), nil
}

func validateTypeName(typeName string) error {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return fmt.Errorf("type name must not be empty")
	}
	if strings.HasPrefix(typeName, "_") || strings.HasSuffix(typeName, "_") {
		return fmt.Errorf("type name must not start or end with an underscore")
	}
	typeNameRegex := regexp.MustCompile(`^[a-z0-9]+(?:_[a-z0-9]+)*$`)
	if !typeNameRegex.MatchString(typeName) {
		return fmt.Errorf("type name must only contain lowercase alphanumeric characters and single underscores (no special characters or spaces)")
	}
	return nil
}
