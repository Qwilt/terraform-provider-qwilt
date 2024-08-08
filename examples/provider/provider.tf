# Configuration resources are implicitly defining dependencies by using references.
# The implicit order is as follows:
# 1. qwiltcdn_certificate, qwiltcdn_site
# 2. qwiltcdn_site_configuration (implicitly depends-on qwiltcdn_site - by using references)
# 3. qwiltcdn_site_activation (implicitly depends-on qwiltcdn_site_configuration, qwiltcdn_certificate - by using references)


provider "qwiltcdn" {
  # Specify username, or set env variable QCDN_USERNAME
  username = "me@mycompany.com"
  # Specify password or set env variable QCDN_PASSWORD
  password = "me123456"
  # Specify token or set env variable QCDN_XAPI_TOKEN
  xapi_token = "eyJhbGciOiJSUzI1NiIsIn..."
}


