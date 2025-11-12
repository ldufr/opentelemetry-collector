package pdata

import (
	"go.opentelemetry.io/collector/internal/cmd/pdatagen/internal/proto"
	"go.opentelemetry.io/collector/internal/cmd/pdatagen/internal/template"
)

const privateByteSliceAccessorTemplate = `// {{ .fieldName }} returns the {{ .fieldName }} associated with this {{ .structName }}.
func (ms {{ .structName }}) {{ .fieldName }}() {{ .packageName }}{{ .returnType }} {
	{{- if .elementHasWrapper }}
	return {{ .packageName }}{{ .returnType }}(internal.New{{ .returnType }}Wrapper(&ms.{{ .origAccessor }}.{{ .originFieldName }}, ms.{{ .stateAccessor }}))
	{{- else }}
	return new{{ .returnType }}(&ms.{{ .origAccessor }}.{{ .originFieldName }}, ms.{{ .stateAccessor }})
	{{- end }}
}`

const privateByteSliceAccessorsTestTemplate = `func Test{{ .structName }}_{{ .fieldName }}(t *testing.T) {
	ms := New{{ .structName }}()
	assert.Equal(t, {{ .packageName }}New{{ .returnType }}(), ms.{{ .fieldName }}())
	ms.{{ .origAccessor }}.{{ .originFieldName }} = internal.GenTest{{ .elementOriginName }}{{ if .elementNullable }}Ptr{{ end }}Slice()
	{{- if .elementHasWrapper }}
	assert.Equal(t, {{ .packageName }}{{ .returnType }}(internal.GenTest{{ .returnType }}Wrapper()), ms.{{ .fieldName }}())
	{{- else }}
	assert.Equal(t, generateTest{{ .returnType }}(), ms.{{ .fieldName }}())
	{{- end }}
}`

const privateByteSliceSetTestTemplate = `orig.{{ .originFieldName }} = internal.GenTest{{ .elementOriginName }}{{ if .elementNullable }}Ptr{{ end }}Slice()`

type PrivateByteSlice struct {
	fieldName string
	protoID   uint32
}

func (pbs *PrivateByteSlice) GenerateAccessors(ms *messageStruct) string {
	t := template.Parse("privateByteSliceAccessorTemplate", []byte(privateByteSliceAccessorTemplate))
	return template.Execute(t, pbs.templateFields(ms))
}

func (pbs *PrivateByteSlice) GenerateAccessorsTest(ms *messageStruct) string {
	t := template.Parse("privateByteSliceAccessorsTestTemplate", []byte(privateByteSliceAccessorsTestTemplate))
	return template.Execute(t, pbs.templateFields(ms))
}

func (pbs *PrivateByteSlice) GenerateTestValue(ms *messageStruct) string {
	t := template.Parse("privateByteSliceSetTestTemplate", []byte(privateByteSliceSetTestTemplate))
	return template.Execute(t, pbs.templateFields(ms))
}

func (pbs *PrivateByteSlice) toProtoField(ms *messageStruct) proto.FieldInterface {
	return &proto.Field{
		Type:              proto.TypeBytes,
		ID:                pbs.protoID,
		Name:              pbs.fieldName,
		MessageName:       byteSlice.getElementOriginName(),
		ParentMessageName: ms.protoName,
		Repeated:          false,
		Nullable:          false,
	}
}

func (pbs *PrivateByteSlice) templateFields(ms *messageStruct) map[string]any {
	return map[string]any{
		"structName":        ms.getName(),
		"fieldName":         pbs.fieldName,
		"originFieldName":   pbs.fieldName,
		"elementOriginName": byteSlice.getElementOriginName(),
		"packageName": func() string {
			if byteSlice.getPackageName() != ms.packageName {
				return byteSlice.getPackageName() + "."
			}
			return ""
		}(),
		"returnType":        byteSlice.getName(),
		"origAccessor":      origAccessor(ms.getHasWrapper()),
		"stateAccessor":     stateAccessor(ms.getHasWrapper()),
		"elementHasWrapper": byteSlice.getHasWrapper(),
		"elementNullable":   byteSlice.getElementNullable(),
	}
}
