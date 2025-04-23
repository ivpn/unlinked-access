## Manifest Generation Service
`generator`

The Manifest Generation Service creates comprehensive, signed lists of all valid tokens along with their associated subscription properties. These manifests serve as the authoritative source of subscription information for mailX and other applications.

Key responsibilities include:

- Creating cryptographically secure tokens for users based on their IVPN user ID
- Managing the secure association between IVPN users and their tokens
- Ensuring that tokens cannot be reversed to reveal user identities
- Supporting token regeneration when necessary for security events
- It incorporates an initial hashing step (SHA512) on the user ID before interacting with the HSM, ensuring the raw user identifier is never sent outside the IVPN server environment to the HSM.
