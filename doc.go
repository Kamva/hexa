// Package hexa is a microservice SDK for the Kamva ecosystem.
//
// It provides the core, transport-agnostic building blocks for building
// microservices: a propagatable request Context (carrying the user,
// correlation id, locale, logger, translator and a concurrency-safe store),
// a structured Error/Reply model, users and user propagation, health
// reporting, a service registry with ordered boot/shutdown, distributed
// locks, and supporting utilities. Most capabilities are defined as
// interfaces with swappable driver implementations in sub-packages.
package hexa
