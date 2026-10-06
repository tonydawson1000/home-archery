// Package postgres will persist Session and Archer via pgx.
//
// Deferred until the application ports and in-memory contract tests are the
// source of truth; this adapter must satisfy the same SessionRepository and
// ArcherRepository interfaces (see application/memory).
package postgres
