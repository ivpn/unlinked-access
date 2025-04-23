## Token Verification Service
`verifier`

The Token Verification Library is integrated into mailX to verify tokens and manage user subscription states based on the information in the manifest. It handles the complexities of token validation without requiring direct communication with IVPN for routine operations.

Key responsibilities:

- Retrieving and validating manifests from the distribution system
- Verifying user tokens against the current manifest
- Managing local subscription state based on manifest data
- Implementing appropriate fallback mechanisms for system failures
- Enforcing grace periods for subscription expiry during outages
- Maintaining access for users with renewed subscriptions during outages
