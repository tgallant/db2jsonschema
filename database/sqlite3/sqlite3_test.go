package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTableSQLSimple(t *testing.T) {
	exampleTable := "CREATE TABLE Example (id int, name varchar(255))"
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 2, len(table.Fields), "the table should have 2 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
}

func TestParseTableSQLSimpleIfNotExists(t *testing.T) {
	exampleTable := `CREATE TABLE IF NOT EXISTS "Example" (id int, name varchar(255))`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 2, len(table.Fields), "the table should have 2 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
}

func TestParseTableSQLWithAttributets(t *testing.T) {
	exampleTable := `
CREATE TABLE Example (
  id int NOT_NULL AUTO_INCREMENT,
  name varchar(255)
)`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 2, len(table.Fields), "the table should have 2 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
}

func TestParseTableSQLWithBackticks(t *testing.T) {
	exampleTable := "CREATE TABLE `Example` (`id` int, `name` varchar(255))"
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 2, len(table.Fields), "the table should have 2 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
}

func TestParseTableSQLWithPrimaryKey(t *testing.T) {
	exampleTable := "CREATE TABLE `Example` (`id` int, `name` varchar(255), PRIMARY KEY (`id`, `name`))"
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 2, len(table.Fields), "the table should have 2 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
	assert.Equal(t, 2, len(table.PrimaryKeys), "there should be 2 primary keys")
	firstPrimaryKey := table.PrimaryKeys[0]
	assert.Equal(t, "id", firstPrimaryKey, "the first primary key should be `id`")
	secondPrimaryKey := table.PrimaryKeys[1]
	assert.Equal(t, "name", secondPrimaryKey, "the second primary key should be `name`")
}

func TestParseTableSQLWithDatetime(t *testing.T) {
	exampleTable := `
CREATE TABLE Example (
  id int NOT_NULL AUTO_INCREMENT,
  name varchar(255),
  created_at datetime
)`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 3, len(table.Fields), "the table should have 3 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
	thirdField := table.Fields[2]
	assert.Equal(t, "created_at", thirdField.Name, "the field name should be `created_at`")
	assert.Equal(t, "string", thirdField.Type.Name, "the field type should be `string`")
	assert.Equal(t, "date-time", thirdField.Type.Format, "the field format should be `date-time`")
}

func TestParseTableSQLWithConstraint(t *testing.T) {
	exampleTable := `
CREATE TABLE Example (
  id int NOT_NULL AUTO_INCREMENT,
  name varchar(255),
  user_id integer,
  team_id integer,
  created_at datetime,
  PRIMARY KEY (id),
  FOREIGN KEY("Test") REFERENCES "Test" (id),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user(id),
  CONSTRAINT fk_team_id FOREIGN KEY (team_id) REFERENCES team(id),
  CHECK ("isDeleted" IN (0, 1))
)`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Example", table.Name, "the table name should be `Example`")
	assert.Equal(t, 5, len(table.Fields), "the table should have 5 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
	secondField := table.Fields[1]
	assert.Equal(t, "name", secondField.Name, "the field name should be `name`")
	assert.Equal(t, "string", secondField.Type.Name, "the field type should be `string`")
	thirdField := table.Fields[2]
	assert.Equal(t, "user_id", thirdField.Name, "the field name should be `user_id`")
	assert.Equal(t, "number", thirdField.Type.Name, "the field type should be `number`")
	fourthField := table.Fields[3]
	assert.Equal(t, "team_id", fourthField.Name, "the field name should be `team_id`")
	assert.Equal(t, "number", fourthField.Type.Name, "the field type should be `number`")
	fifthField := table.Fields[4]
	assert.Equal(t, "created_at", fifthField.Name, "the field name should be `created_at`")
	assert.Equal(t, "string", fifthField.Type.Name, "the field type should be `string`")
	assert.Equal(t, "date-time", fifthField.Type.Format, "the field format should be `date-time`")
}

func TestParseTableSQLFromSQLAlchemy(t *testing.T) {
	exampleTable := `
CREATE TABLE IF NOT EXISTS "Condition" (
  id INTEGER NOT_NULL,
  kind VARCHAR(255) NOT NULL,
  "conditionType" VARCHAR NOT NULL,
  "WorkflowStepProgressionId" INTEGER,
  "isDeleted" BOOLEAN,
  "deletedAt" DATETIME,
  "createdAt" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  "updatedAt" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY("WorkflowStepProgressionId") REFERENCES "WorkflowStepProgression" (id),
  CONSTRAINT fk_team_id PRIMARY KEY (team_id),
  CONSTRAINT fk_team_id PRIMARY KEY (team_id) REFERENCES team(id),
  CHECK ("isDeleted" IN (0, 1))
)`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "Condition", table.Name, "the table name should be `Example`")
	assert.Equal(t, 8, len(table.Fields), "the table should have 5 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
}

func TestParseTableSQLFromAlembic(t *testing.T) {
	exampleTable := `
CREATE TABLE "WorkflowTemplate" (
  id INTEGER,
  "isDeleted" BOOLEAN,
  "deletedAt" DATETIME,
  CHECK ("isDeleted" IN (0, 1))
)`
	table, err := ParseTableSQL(exampleTable)
	assert.Nil(t, err, "parsing the table sql should succeed")
	assert.Equal(t, "WorkflowTemplate", table.Name, "the table name should be `Example`")
	assert.Equal(t, 3, len(table.Fields), "the table should have 5 fields")
	firstField := table.Fields[0]
	assert.Equal(t, "id", firstField.Name, "the field name should be `id`")
	assert.Equal(t, "number", firstField.Type.Name, "the field type should be `number`")
}
