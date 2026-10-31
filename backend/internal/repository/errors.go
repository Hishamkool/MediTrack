// Package repository contains all SQL. Handlers never write SQL directly.
// Every query uses $1, $2... placeholders (never string concatenation),
// which protects against SQL injection.
package repository

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)
