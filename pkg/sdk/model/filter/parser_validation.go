// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"fmt"
	"strings"

	"github.com/banzaicloud/logging-operator/pkg/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var validParserTypes = sets.NewString(
	"apache2", "apache_error", "nginx", "syslog", "csv", "tsv", "ltsv", "json", "multiline", "none", "logfmt", "multi_format", "regexp",
)

var validSingleParserTypes = sets.NewString(
	"apache2", "apache_error", "nginx", "syslog", "csv", "tsv", "json", "none", "logfmt", "regexp",
)

var validParserFieldTypes = sets.NewString(
	"string", "integer", "float", "bool", "time", "array",
)

var validParserTimeTypes = sets.NewString(
	"float", "unixtime", "string", "",
)

// ValidateFields performs validation on each field of the object
func (c *ParserConfig) ValidateFields(path *field.Path) field.ErrorList {
	// Empty config is valid
	if c == nil {
		return nil
	}

	// Return any errors
	return validateParserConfig(path, *c)
}

func validateParserConfig(path *field.Path, parser ParserConfig) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Validate KeyName
	{
		// Unvalidated - can be any value
	}

	// Validate ReserveTime
	{
		// Unvalidated - simple boolean
	}

	// Validate ReserveData
	{
		// Unvalidated - simple boolean
	}

	// Validate RemoveKeyNameField
	{
		// Unvalidated - simple boolean
	}

	// Validate ReplaceInvalidSequence
	{
		// Unvalidated - simple boolean
	}

	// Validate InjectKeyPrefix
	{
		// Unvalidated - can be any value
	}

	// Validate HashValueField
	{
		// Unvalidated - can be any value
	}

	// Validate EmitInvalidRecordToError
	{
		// Unvalidated - simple boolean
	}

	// Validate Parse
	{
		allErrs = append(allErrs, validateParseSection(path.Child("parse"), parser.Parse)...)
	}

	// Validate Parsers
	{
		// This field is deprecated, error on its use
		if len(parser.Parsers) > 0 {
			allErrs = append(allErrs, field.Invalid(path.Child("parses"), parser.Parsers, "'parsers' field is deprecated, use 'parse'"))
		}
	}

	// Return any errors
	return allErrs
}

func validateParseSection(path *field.Path, section ParseSection) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Validate Type
	{
		if !validParserTypes.Has(section.Type) {
			allErrs = append(allErrs, field.NotSupported(path.Child("type"), section.Type, validParserTypes.List()))
		}
	}

	// Validate Expression
	{
		if section.Expression == "" && section.Type == "regex" {
			allErrs = append(allErrs, field.Required(path.Child("expression"), "'expression' is required for regexp parsers"))
		}

		if section.Expression != "" && section.Type != "regex" {
			allErrs = append(allErrs, field.Required(path.Child("expression"), "'expression' can only be specified for regexp parsers"))
		}

		if err := validation.ValidateFluentdRegex(section.Expression); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("expression"), section.Expression, err.Error()))
		}
	}

	// Validate TimeKey
	{
		// Unvalidated - can be any value
	}

	// Validate NullValuePattern
	{
		if err := validation.ValidateFluentdRegex(section.NullValuePattern); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("null_value_pattern"), section.NullValuePattern, err.Error()))
		}
	}

	// Validate NullEmptyString
	{
		// Unvalidated - simple boolean
	}

	// Validate EstimateCurrentEvent
	{
		// Unvalidated - simple boolean
	}

	// Validate KeepTimeKey
	{
		// Unvalidated - simple boolean
	}

	// Validate Types
	{
		// Parse the hash
		types, err := validation.ParseFluentdHash(section.Types)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("types"), section.Types, err.Error()))
		}

		// Validate the types, this code is based on the fluentd equivalent
		// https://github.com/fluent/fluentd/blob/be2d977b917b7a109bb02d0b2f2503f81fff16fc/lib/fluent/plugin/parser.rb#L229
		for name, typ := range types {
			// Split the value on ":" to extract information
			parts := strings.SplitN(typ, ":", 2)
			typeName := parts[0]

			// Ensure
			if !validParserFieldTypes.Has(typeName) {
				msg := fmt.Sprintf("field %q using unsupported type %q, valid types include %v", name, typeName, validParserFieldTypes.List())
				allErrs = append(allErrs, field.Invalid(path.Child("types"), section.Types, msg))
				continue
			}

			// Time fields have an extra option that needs validating
			if typeName == "time" {
				// TODO
			}
		}
	}

	// Validate TimeFormat
	{
		// TODO: this can be multiple different formats, investigate how fluentd handles this
	}

	// Validate TimeType
	{
		if !validParserTimeTypes.Has(section.TimeType) {
			allErrs = append(allErrs, field.NotSupported(path.Child("time_type"), section.TimeType, validParserTimeTypes.List()))
		}
	}

	// Validate LocalTime
	{
		if section.LocalTime && section.UTC {
			allErrs = append(allErrs, field.Invalid(path.Child("local_time"), section.LocalTime, "'local_time' and 'utc' cannot both be true"))
		}
	}

	// Validate UTC
	{
		if section.LocalTime && section.UTC {
			allErrs = append(allErrs, field.Invalid(path.Child("utc"), section.UTC, "'local_time' and 'utc' cannot both be true"))
		}
	}

	// Validate Timezone
	{
		// TODO: Investigate how fluentd parses this
	}

	// Validate Format
	{
		// Format is only a valid field in the patterns, not in the <parse> block
		if section.Format != "" {
			allErrs = append(allErrs, field.Invalid(path.Child("format"), section.Format, "format can only be specified in patterns"))
		}
	}

	// Validate FormatFirstline
	{
		// This field is part of the multiline plugin, so only valid when the type is multiline
		if section.FormatFirstline != "" && section.Type != "multiline" {
			allErrs = append(allErrs, field.Invalid(path.Child("format_firstline"), section.FormatFirstline, "'format_firstline' can only be specified for multiline parsers"))
		}

		// Validate the regex
		if err := validation.ValidateFluentdRegex(section.FormatFirstline); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("format_firstline"), section.FormatFirstline, err.Error()))
		}
	}

	// Validate Delimiter
	{
		// Only available when using type: ltsv (default: "\t")
		if section.Delimiter != "" && section.Type != "ltsv" {
			allErrs = append(allErrs, field.Invalid(path.Child("delimiter"), section.Delimiter, "'delimiter' can only be specified for ltsv parsers"))
		}

		// Field is redundant when DelimiterPattern is set
		if section.Delimiter != "" && section.DelimiterPattern != "" {
			allErrs = append(allErrs, field.Invalid(path.Child("delimiter"), section.Delimiter, "'delimiter' can not be specified when a 'delimiter_pattern' is provided"))
		}
	}

	// Validate DelimiterPattern
	{
		// Only available when using type: ltsv
		if section.DelimiterPattern != "" && section.Type != "ltsv" {
			allErrs = append(allErrs, field.Invalid(path.Child("delimiter_pattern"), section.DelimiterPattern, "'delimiter_pattern' can only be specified for ltsv parsers"))
		}

		// Validate the regex
		if err := validation.ValidateFluentdRegex(section.DelimiterPattern); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("delimiter_pattern"), section.DelimiterPattern, err.Error()))
		}
	}

	// Validate LabelDelimiter
	{
		// Only available when using type: ltsv
		if section.LabelDelimiter != "" && section.Type != "ltsv" {
			allErrs = append(allErrs, field.Invalid(path.Child("label_delimiter"), section.LabelDelimiter, "'label_delimiter' can only be specified for ltsv parsers"))
		}
	}

	// Validate Multiline
	{
		// Only available when using type: multiline
		if len(section.Multiline) > 0 && section.Type != "multiline" {
			allErrs = append(allErrs, field.Invalid(path.Child("multiline"), section.Multiline, "'multiline' can only be specified for multiline parsers"))
		}

		// Validate each lines regex
		for i, line := range section.Multiline {
			if err := validation.ValidateFluentdRegex(line); err != nil {
				allErrs = append(allErrs, field.Invalid(path.Child("multiline").Index(i), line, err.Error()))
			}
		}
	}

	// Validate Patterns
	{
		for i, pattern := range section.Patterns {
			allErrs = append(allErrs, validateSingleParseSection(path.Child("patterns").Index(i), pattern)...)
		}
	}

	// Return any errors
	return allErrs
}

func validateSingleParseSection(path *field.Path, section SingleParseSection) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Validate Type
	{
		if section.Type != "" {
			allErrs = append(allErrs, field.Invalid(path.Child("type"), section.Format, "'type' cannot be specified in patterns, use 'format'"))
		}
	}

	// Validate Expression
	{
		if section.Expression == "" && section.Format == "regex" {
			allErrs = append(allErrs, field.Required(path.Child("expression"), "'expression' is required for regexp parsers"))
		}

		if section.Expression != "" && section.Format != "regex" {
			allErrs = append(allErrs, field.Required(path.Child("expression"), "'expression' can only be specified for regexp parsers"))
		}

		if err := validation.ValidateFluentdRegex(section.Expression); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("expression"), section.Expression, err.Error()))
		}
	}

	// Validate TimeKey
	{
		// Unvalidated - can be any value
	}

	// Validate NullValuePattern
	{
		if err := validation.ValidateFluentdRegex(section.NullValuePattern); err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("null_value_pattern"), section.NullValuePattern, err.Error()))
		}
	}

	// Validate NullEmptyString
	{
		// Unvalidated - simple boolean
	}

	// Validate EstimateCurrentEvent
	{
		// Unvalidated - simple boolean
	}

	// Validate KeepTimeKey
	{
		// Unvalidated - simple boolean
	}

	// Validate Types
	{
		// Parse the hash
		types, err := validation.ParseFluentdHash(section.Types)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("types"), section.Types, err.Error()))
		}

		// Validate the types, this code is based on the fluentd equivalent
		// https://github.com/fluent/fluentd/blob/be2d977b917b7a109bb02d0b2f2503f81fff16fc/lib/fluent/plugin/parser.rb#L229
		for name, typ := range types {
			// Split the value on ":" to extract information
			parts := strings.SplitN(typ, ":", 2)
			typeName := parts[0]

			// Ensure
			if !validParserFieldTypes.Has(typeName) {
				msg := fmt.Sprintf("field %q using unsupported type %q, valid types include %v", name, typeName, validParserFieldTypes.List())
				allErrs = append(allErrs, field.Invalid(path.Child("types"), section.Types, msg))
				continue
			}

			// Time fields have an extra option that needs validating
			if typeName == "time" {
				// TODO
			}
		}
	}

	// Validate TimeFormat
	{
		// TODO: this can be multiple different formats, investigate how fluentd handles this
	}

	// Validate TimeType
	{
		if !validParserTimeTypes.Has(section.TimeType) {
			allErrs = append(allErrs, field.NotSupported(path.Child("time_type"), section.TimeType, validParserTimeTypes.List()))
		}
	}

	// Validate LocalTime
	{
		if section.LocalTime && section.UTC {
			allErrs = append(allErrs, field.Invalid(path.Child("local_time"), section.LocalTime, "'local_time' and 'utc' cannot both be true"))
		}
	}

	// Validate UTC
	{
		if section.LocalTime && section.UTC {
			allErrs = append(allErrs, field.Invalid(path.Child("utc"), section.UTC, "'local_time' and 'utc' cannot both be true"))
		}
	}

	// Validate Timezone
	{
		// TODO: Investigate how fluentd parses this
	}

	// Validate Format
	{
		if !validSingleParserTypes.Has(section.Format) {
			allErrs = append(allErrs, field.NotSupported(path.Child("format"), section.Format, validSingleParserTypes.List()))
		}
	}

	// Return any errors
	return allErrs
}
