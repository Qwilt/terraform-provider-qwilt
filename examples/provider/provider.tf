#The Qwilt Terraform Provider supports two authentication methods:

#- **API key-based authentication** - The preferred method.
#    - When the *api_key* parameter is set, the key is passed 
#      in the header of each API call to authenticate the request. 
#    - To obtain an API key, please contact support@qwilt.com. 

#- **Login with username and password** - Not recommended. 
#  -  When the *user name* and *password* parameters are set, any 
#     Terraform command (apply, refresh, plan, etc.)  triggers the
#     Qwilt Login API to generate the required cqloud access token. 
#  -  Support for this method may be deprecated in the future.


# You can set the authentication parameters inside the provider 
# configuration or as environment variables. 
# We recommend setting env variables.

# |TF Provider Variable |  Env Variable   | Example Value |
# | --- | --- | --- |
# |api_key | QCDN_API_KEY |  "eyJhbGciOiJSUzI1NiIsIn..." |
# | username| QCDN_USERNAME  | "me@mycompany.com" |
# | password |QCDN_PASSWORD |  "mypwd123456" |


#**Notes**:
#- If the QCDN_API_KEY env variable is defined, the QCDN_USERNAME 
#  and QCDN_PASSWORD env variables are ignored. 
#- If you set the authentication parameters in the Terraform provider configuration,
#  you can define *either* the api_key *or* the username and password. 

# Example of how to set the QCDN_API_KEY env variable:
#
#     export QCDN_API_KEY="eyJhbGciOiJSUzI1NiIsIn..."
#
#
# When the authentication parameters are set by the environment variables, the provider config looks like this:
#  
#     provider "qwilt" { }
#  


# Example of how to set the API Key param in the provider config:
#
#    provider "qwilt" {
#       api_key = "eyJhbGciOiJSUzI1NiIsIn..."
#    }


# Example of how to set the username and password params in the provider config:
#  
#    provider "qwilt" {
#        username = "me@mycompany.com"
#        password = "me123456"
#    }
# 


