# authentication is supported by methods: login and xapi token
# 1. api_key is always the preferred method
# 2. user/password are supported but not recommended and might be deprecated in the near future.
# 3. each of these variables can be replaced with a corresponding env variable:
#  - username --> QCDN_USERNAME
#  - password --> QCDN_PASSWORD
#  - api_key --> QCDN_API_KEY

provider "qwilt" {
  # Specify username, or set env variable QCDN_USERNAME
  username = "me@mycompany.com"
  # Specify password or set env variable QCDN_PASSWORD
  password = "me123456"
  # Specify token or set env variable QCDN_API_KEY
  api_key = "eyJhbGciOiJSUzI1NiIsIn..."
}


