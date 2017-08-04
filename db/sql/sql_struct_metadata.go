// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import "reflect"

// TableName interface is used to specify custom table name for SQL statements
type TableName interface {
	// TableName returns sql table name.
	TableName() string
}

// SchemaName interface is used to specify custom schema name for SQL statements
type SchemaName interface {
	// SchemaName returns sql schema name where the table resides
	SchemaName()  string
}

// EntityTableName tries to cast to SchemaName & TableName interfaces.
// If not possible it uses just name of struct as table name.
func EntityTableName(entity interface{}) string {
	var tableName, schemaName string
	if nameProvider, ok := entity.(TableName); ok {
		tableName = nameProvider.TableName()
	}

	if tableName == "" {
		tableName = reflect.Indirect(reflect.ValueOf(entity)).Type().Name()
	}

	if schemaNameProvider, ok := entity.(SchemaName); ok {
		schemaName = schemaNameProvider.SchemaName()
	}

	if schemaName == "" {
		return tableName
	}

	return  schemaName + "." + tableName
}