#Create an empty resource to import into.

#After the import is complete, manually set the required attributes 
#in the resource based on the imported state.

#We recommend changing the site_id and revision_id attributes to references
#to the qwilt_cdn_site_configuration resource to achieve implicit dependency.


resource "qwilt_cdn_site_activation" "example" {
}

# Import the qwilt_cdn_site_activation resource by specifying the site_id. 

# For example: terraform import qwilt_cdn_site_activation.example <site_id>

    # The process determines which configuration to import based on
    # the following conditions: 
    # - If there is an active published site configuration, it is imported.
    # - If not, the most recently saved configuration version is imported.
        
# Alternatively, you can specify a particular publish_id by adding 
# a colon (:) and the publish_id. 
    
# For example: 
#    */terraform import qwilt_cdn_site_activation.example xxxxxxxxxxxx:yyyyyyyyyyyy

terraform import qwilt_cdn_site_activation.example <site_id>:<publish_id>

