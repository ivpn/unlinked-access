## Token Generation Service
`token`

The Token Generation Service is responsible for creating secure, non-linkable tokens that allow IVPN users to authenticate with mailX without revealing their identity. This component interfaces directly with the Hardware Security Module (HSM) to ensure that the cryptographic operations remain secure and tamper-resistant.

Key responsibilities:

- Creating cryptographically secure tokens for users based on their IVPN user ID
- Managing the secure association between IVPN users and their tokens
- Ensuring that tokens cannot be reversed to reveal user identities
- Supporting token regeneration when necessary for security events
- It incorporates an initial hashing step (SHA512) on the user ID before interacting with the HSM, ensuring the raw user identifier is never sent outside the IVPN server environment to the HSM.
