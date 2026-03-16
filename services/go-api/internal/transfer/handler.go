// Purpose: HTTP handlers for Cross-Position Transfer — initiate transfer, validate target position,
//
//	carry over history chain, block test sends from closed source application.
//
// Ref: API Contract — POST /api/v1/applications/{id}/transfer
// DB tables: applications, application_transfers
// Dev must implement: InitiateTransfer, ValidateTarget, GetTransferChain
package transfer
