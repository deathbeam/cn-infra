//  Copyright (c) 2018 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package agent

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"go.ligato.io/cn-infra/v2/infra"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/cn-infra/v2/logging/logrus"
)

var infraLogger = logrus.NewLogger("infra")

func init() {
	if os.Getenv("DEBUG_INFRA") != "" {
		infraLogger.SetLevel(logging.DebugLevel)
		infraLogger.Debugf("infra debug logger enabled")
	}
}

var (
	// use DEBUG_INFRA=lookup to print plugin verbose lookup logs
	printPluginLookupDebugs = strings.Contains(strings.ToLower(os.Getenv("DEBUG_INFRA")), "lookup")
	// use DEBUG_INFRA=start to print plugin start durations
	printPluginStartDurations = strings.Contains(strings.ToLower(os.Getenv("DEBUG_INFRA")), "start")
)

func findPlugins(val reflect.Value, uniqueness map[infra.Plugin]struct{}, x ...int) (
	res []infra.Plugin, err error,
) {
	n := 0
	if len(x) > 0 {
		n = x[0]
	}
	var logf = func(f string, a ...interface{}) {
		for i := 0; i < n; i++ {
			f = "\t" + f
		}
		if printPluginLookupDebugs {
			infraLogger.Debugf(f, a...)
		}
	}

	typ := val.Type()

	logf("=> %v (%v)", typ, typ.Kind())
	defer logf("== %v ", typ)

	if typ.Kind() == reflect.Interface {
		if val.IsNil() {
			logf(" - val is nil")
			return nil, nil
		}
		val = val.Elem()
		typ = val.Type()
		//logf(" - interface to elem: %v (%v)", typ, val.Kind())
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		//logrus.DefaultLogger().Debug(" - typ ptr kind: ", typ)
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		//logrus.DefaultLogger().Debug(" - val ptr kind: ", val)
	}

	if !val.IsValid() {
		logf(" - val is invalid")
		return nil, nil
	}

	if typ.Kind() != reflect.Struct {
		logf(" - is not a struct: %v %v", typ.Kind(), val.Kind())
		return nil, nil
	}

	//logf(" -> checking %d fields", typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// PkgPath is empty for exported fields
		if exported := field.PkgPath == ""; !exported {
			continue
		}

		fieldVal := val.Field(i)

		logf("-> field %d: %v - %v (%v)", i, field.Name, field.Type, fieldVal.Kind())

		// transform field to list of values if it is slice or array
		fieldVals, isList := getFieldValues(field, fieldVal)
		if isList {
			logf(" - found list: %v", field.Name)
		}

		for _, entry := range fieldVals {
			var fieldPlug infra.Plugin
			plug, implementsPlugin := isFieldPlugin(entry.fieldVal)
			if implementsPlugin {
				if plug == nil {
					logf(" - found nil plugin: %v", entry.fieldName)
					continue
				}

				_, found := uniqueness[plug]
				if found {
					logf(" - found duplicate plugin: %v %v", entry.fieldName, field.Type)
					continue
				}

				// TODO: perhaps add regexp for validation of plugin name

				uniqueness[plug] = struct{}{}
				fieldPlug = plug

				logf(" + FOUND PLUGIN: %v - %v (%v)", plug.String(), entry.fieldName, field.Type)
			}

			// do recursive inspection only for plugins and fields Deps
			if fieldPlug != nil || (field.Anonymous && entry.fieldVal.Kind() == reflect.Struct) {
				// try to inspect structure recursively
				l, err := findPlugins(entry.fieldVal, uniqueness, n+1)
				if err != nil {
					logf(" - Bad field: %v %v", entry.fieldName, err)
					continue
				}
				//logf(" - listed %v plugins from %v (%v)", len(l), field.Name, field.Type)
				res = append(res, l...)
			}

			if fieldPlug != nil {
				res = append(res, fieldPlug)
			}
		}
	}

	logf("<- got %d plugins", len(res))

	return res, nil
}

type fieldValEntry struct {
	fieldName string
	fieldVal  reflect.Value
}

func getFieldValues(field reflect.StructField, fieldVal reflect.Value) ([]*fieldValEntry, bool) {
	var fieldVals []*fieldValEntry

	kind := fieldVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < fieldVal.Len(); i++ {
			entryName := fmt.Sprintf("%s[%d]", field.Name, i)
			fieldVals = append(fieldVals, &fieldValEntry{entryName, fieldVal.Index(i)})
		}
		return fieldVals, true
	}

	// underlying list
	typ := reflect.ValueOf(fieldVal.Interface()).Kind()
	if typ == reflect.Slice || typ == reflect.Array {
		fieldValSlice := reflect.ValueOf(fieldVal.Interface())
		for i := 0; i < fieldValSlice.Len(); i++ {
			entryName := fmt.Sprintf("%s[%d]", field.Name, i)
			fieldVals = append(fieldVals, &fieldValEntry{entryName, fieldValSlice.Index(i)})
		}
		return fieldVals, true
	}

	// not a list (single entry)
	fieldVals = append(fieldVals, &fieldValEntry{field.Name, fieldVal})
	return fieldVals, false
}

var pluginType = reflect.TypeOf((*infra.Plugin)(nil)).Elem()

func isFieldPlugin(fieldVal reflect.Value) (infra.Plugin, bool) {
	//logrus.DefaultLogger().Debugf(" - is field plugin: %v (%v) %v", field.Type, fieldVal.Kind(), fieldVal)

	switch fieldVal.Kind() {
	case reflect.Struct:
		ptrType := reflect.PtrTo(fieldVal.Type())
		if ptrType.Implements(pluginType) {
			if fieldVal.CanAddr() {
				if plug, ok := fieldVal.Addr().Interface().(infra.Plugin); ok {
					return plug, true
				}
			}
			return nil, true
		}
	case reflect.Ptr, reflect.Interface:
		if plug, ok := fieldVal.Interface().(infra.Plugin); ok {
			if fieldVal.IsNil() {
				return nil, true
			}
			return plug, true
		}
	}

	return nil, false
}
