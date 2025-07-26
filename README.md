# Auth System
Private Cross-Application Authentication System (PCAAS)

## Services

### Token Generation Service (server)
`token`

The Token Generation Service is responsible for creating secure, non-linkable tokens that allow IVPN users to authenticate with MailX without revealing their identity. This component interfaces directly with the Hardware Security Module (HSM) to ensure that the cryptographic operations remain secure and tamper-resistant.

### Manifest Generation Service (server)
`generator`

The Manifest Generation Service creates comprehensive, signed lists of all valid tokens along with their associated subscription properties. These manifests serve as the authoritative source of subscription information for MailX and other applications.

### Manifest Distribution Service (server)
`distributor`

The Manifest Distribution System securely delivers manifests to authorized applications. It ensures that only legitimate applications can access the manifest while optimizing distribution for performance and reliability.

### Pre-Authorization Service (server)
`preauth`

The Pre-Authorization API provides immediate verification for new user signups, ensuring that users can access MailX instantly after signup rather than waiting for the next manifest update.

### Token Verification Service (client)
`verifier`

The Token Verification Library is integrated into MailX to verify tokens and manage user subscription states based on the information in the manifest. It handles the complexities of token validation without requiring direct communication with IVPN for routine operations.

## Installation

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Usage
1. Clone the repo: `git clone <repo-url> && cd <repo-dir>`
2. Create and configure `.env`: `cp .env.sample .env`
3. Start services: `docker compose up -d`
4. View logs: `docker compose logs -f`
5. Stop services: `docker compose down`
