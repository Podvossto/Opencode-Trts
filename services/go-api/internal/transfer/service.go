// Purpose: Business logic for Transfer — close source application, return headcount slot,
//
//	create target application with prev_application_id chain, enforce interview-stage-only rule.
//
// Ref: US-M09-01 through US-M09-04
// Dev must implement: TransferService interface — Transfer(sourceAppID, targetJobID uuid.UUID) error
package transfer
