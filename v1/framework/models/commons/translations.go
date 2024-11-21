package commons

import (
	pbte "github.com/Allen-Career-Institute/common-protos/resource/v1/types/enums"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

const (
	AcadSessionSeperator = "-"
	YearFormat           = "yyyy"
	YearShortFormat      = "yy"
)

func GetEnumDisplayName(enum interface{ protoreflect.Enum }) string {
	displayName := proto.GetExtension(enum.Descriptor().Values().ByNumber(enum.Number()).Options(), pbte.E_DisplayName)
	return displayName.(string)
}

func GetEnumExtension[T any](enum interface{ protoreflect.Enum }, xt protoreflect.ExtensionType) T {
	displayName := proto.GetExtension(enum.Descriptor().Values().ByNumber(enum.Number()).Options(), xt)
	return displayName.(T)
}

func GetAcadSessionShort(session string) string {
	var response []string
	sessions := strings.Split(session, AcadSessionSeperator)
	for _, sess := range sessions {
		if len(sess) > 2 {
			response = append(response, sess[2:])
		}
	}
	return strings.Join(response, AcadSessionSeperator)
}
