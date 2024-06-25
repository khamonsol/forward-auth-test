## Project Description

### Overview
The project aims to build a standalone forward authentication server that integrates with Traefik to validate user permissions based on Azure tokens using OIDC. The server will utilize externalized policy and authentication provider configurations, allowing for flexible and extensible authentication mechanisms.

#### Summary
This project will deliver a forward authentication server that:
- Validates user permissions using OIDC with Azure tokens.
- Uses externalized policy and authentication provider configurations.
- Provides detailed logging and tracking with UUID correlation IDs.
- Is extensible for future authentication providers and validation steps.


### Key Components
1. **Auth Provider Configs**: Configurations aligned with Azure enterprise applications, modeled as a factory for easy implementation of new providers.
2. **Policy Management**: Policies that determine user access based on roles and users, retrieved from Kubernetes ConfigMaps.
3. **Request Handling**: HTTP handlers that validate tokens, enforce policies, and maintain correlation IDs for tracking requests.

## Requirements

### Functional Requirements
1. **Auth Provider Configs**:
    - Implement a factory pattern to support multiple authentication providers.
    - Providers must load configuration from Kubernetes ConfigMaps.
    - Providers should include methods to get the provider name and issuer URL.

2. **Policy Management**:
    - Policies define the roles and users that can access specific endpoints.
    - Policies are retrieved based on the host and HTTP verb, using a deterministic string for the ConfigMap key.
    - Policies are stored in Kubernetes ConfigMaps and should be loaded and parsed correctly.

3. **Request Handling**:
    - Implement HTTP handler functions for each validation step to allow for future extensibility.
    - Discover and validate auth tokens from headers (`Authorization`, `x-access-token`, `X-Access-Token`).
    - Maintain the token on the request if the request is allowed.

### Ancillary Requirements
1. **UUID Correlation ID**:
    - Add a UUID correlation ID to all requests passed into the forward auth.
    - Return the correlation ID as an identifier for tracking and add it to logs.
    - Correlation ID should be added to the logs in the `correlation_id` field.

2. **JSON Logging**:
    - All logs should be in JSON format using the slog library.

3. **Header Handling**:
    - On the request sent to the backend (if the request is allowed), maintain the token on whichever header it was passed in on.


## Usage

Minimally the forward auth server will need an auth provider configured along with an access policy. Below are examples

**NOTE:** Config maps are expected to be present in the same namespace as the forward auth server. 

Sample provider config: Note client_secret is not yet implmented, if we need it, we will implement it as a secret that
is pulled from vault. 


**auth provider config map:** azure_beyond_prod_config
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: azure_beyond_prod_config
data:
  client_id: "your-actual-client-id"
  tenant_id: "your-actual-tenant-id"
```


Sample config map that allows a user with the role API_PNL_RO access to positions and pnl detail api:

**access policy config map:** access_policy_beyond_soleaenergy_com
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: access_policy_beyond_soleaenergy_com
data:
  api_v1_pos_pnl_detail_post: |
    provider_name: azure_beyond_prod
    provider_type: azure
    roles:
      - PNL_API_RW
    users: []
```

