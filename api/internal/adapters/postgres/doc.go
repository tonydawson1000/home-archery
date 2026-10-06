// Package postgres persists Session and Archer via pgx.
//
// It implements the same ArcherRepository and SessionRepository ports as
// application/memory. Scoring rules stay in domain; SQL stores the aggregate
// and projection columns.
package postgres
