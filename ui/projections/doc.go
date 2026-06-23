// Package projections implements Dogma projection handlers that build the
// read-models used by the UI.
//
// These projections serve the web interface only and do not contain any
// business logic. They are implemented using the SQL adapter from
// [projectionkit], a companion library for Dogma that simplifies building
// projection handlers for various backing stores. The SQL adapter manages
// transaction handling, offset tracking and idempotency, so each handler only
// needs to provide the SQL statements that update the read-model tables.
//
// [projectionkit]: https://github.com/dogmatiq/projectionkit
package projections
