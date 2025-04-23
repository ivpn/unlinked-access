## Manifest Generation Service
`generator`

The Manifest Generation Service creates comprehensive, signed lists of all valid tokens along with their associated subscription properties. These manifests serve as the authoritative source of subscription information for MailX and other applications.

Key responsibilities:

- Collecting current subscription status information for all IVPN users
- Creating a structured manifest document containing token hashes and subscription details
- Signing the manifest using the HSM to ensure authenticity and integrity
- Managing manifest versioning and history
- Scheduling regular manifest updates on an hourly basis
