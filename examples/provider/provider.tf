# Both API key-based and token-based authentication is supported.
# 1. api_key is always the preferred method
# 2. user/password (token-based) are supported but not recommended and might be deprecated in the near future.
# 3. each of these variables can be replaced with a corresponding env variable:
#  - username --> QCDN_USERNAME
#  - password --> QCDN_PASSWORD
#  - api_key --> QCDN_API_KEY

provider "qwilt" {
  # Specify username and password, or set env variables QCDN_USERNAME and QCDN_PASSWORD.
  # username = "me@mycompany.com"
  # password = "me123456"
  # Or, specify key or set env variable QCDN_API_KEY.
  api_key = "eyJhbGciOiJSUzI1NiIsIn..."
}


#Provider configuration: If you set the authentication parameters in the Terraform provider configuration, 
#you can define either the username and password *or* the api_key.

# Environment variables: If the QCDN_API_KEY env variable is defined, then the QCDN_USERNAME and QCDN_PASSWORD env variables are ignored.

# When the username and password parameters are set, any Terraform command (apply, refresh, plan, etc.) triggers the Qwilt Login API to 
# generate the required cqloud access token. A token is valid for one hour.

# When the api_key parameter is set, the key is passed in the header of each API call to authenticate the request. A key is valid for one year.
# To obtain an API key, please contact [support@qwilt.com](mailto:support@qwilt.com).

# When the authentication parameters are set by the shell environment variables, the provider config looks like this:
# ```
# provider "qwilt" {}
# ```

